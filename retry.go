// Hand-written, not generated (see overlay/): an opt-in retry layer for the HTTP client, and idempotency keys.
// The behaviour is the same in every OpenRoam SDK; docs/RETRIES.md in the SDK repository describes it.

package openroam

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	mathrand "math/rand/v2"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// RetryOptions configures NewRetryingHTTPClient and RetryTransport. A zero field takes its default.
type RetryOptions struct {
	// MaxRetries is the number of retries after the first attempt for network errors, 429, 502, 503 and 504.
	// Zero means the default, 3 (at most 4 attempts); a negative value means no such retries.
	MaxRetries int
	// BaseDelay is the first backoff step. Zero means 500 ms.
	BaseDelay time.Duration
	// MaxDelay is the upper bound of any wait, including Retry-After. Zero means 30 s.
	MaxDelay time.Duration
	// InProgressRetries is the number of retries while the server answers 12002 (IDEMPOTENCY_KEY_IN_PROGRESS).
	// Zero means the default, 10; a negative value means none.
	InProgressRetries int
	// OnRetry, when set, is called before each wait: attempt is the 1-based number of the retry about to be made,
	// reason says why (for example "status 503"), delay is how long it waits.
	OnRetry func(attempt int, reason string, delay time.Duration)

	// Test hooks.
	sleep  func(ctx context.Context, d time.Duration) error
	random func() float64 // in [0, 1)
	now    func() time.Time
}

const (
	defaultMaxRetries        = 3
	defaultBaseDelay         = 500 * time.Millisecond
	defaultMaxDelay          = 30 * time.Second
	defaultInProgressRetries = 10
	defaultInProgressDelay   = 2 * time.Second
	codeInProgress           = 12002
	// How much of a response body is read to find its error code, and drained before a connection is reused.
	maxInspectedBody = 1 << 20
)

// RetryTransport is an http.RoundTripper that repeats a request when that is safe: always for a refused request
// (429), and for network errors, 502, 503, 504 and 12002 (still in progress) only when the request is a GET or
// HEAD or carries an Idempotency-Key header. It never retries 500 or any other 4xx. After the last allowed
// attempt it returns the last response (or error) unchanged.
type RetryTransport struct {
	// Base sends each attempt. Nil means http.DefaultTransport.
	Base    http.RoundTripper
	Options RetryOptions
}

// NewRetryingHTTPClient returns a copy of base (nil means a new http.Client) whose transport retries as
// RetryTransport does. Use it as the Configuration's HTTPClient.
func NewRetryingHTTPClient(base *http.Client, opts RetryOptions) *http.Client {
	client := &http.Client{}
	if base != nil {
		*client = *base
	}
	client.Transport = &RetryTransport{Base: client.Transport, Options: opts}
	return client
}

// NewIdempotencyKey returns a new random UUID (version 4) to send as the Idempotency-Key of one purchase.
// Make it once, before the first attempt, and keep it for that purchase.
func NewIdempotencyKey() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic("openroam: no random source for an idempotency key: " + err.Error())
	}
	b[6] = b[6]&0x0f | 0x40 // version 4
	b[8] = b[8]&0x3f | 0x80 // RFC 4122 variant
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// RoundTrip sends req, repeating it as described on RetryTransport.
func (t *RetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	o := t.Options.withDefaults()
	base := t.Base
	if base == nil {
		base = http.DefaultTransport
	}
	ctx := req.Context()

	getBody, err := replayableBody(req)
	if err != nil {
		return nil, err
	}
	repeatable := req.Method == http.MethodGet || req.Method == http.MethodHead || req.Header.Get("Idempotency-Key") != ""

	retries, inProgress := 0, 0
	for attempt := 1; ; attempt++ {
		try := req.Clone(ctx)
		if getBody != nil {
			if try.Body, err = getBody(); err != nil {
				return nil, err
			}
			try.GetBody = getBody
		}

		resp, err := base.RoundTrip(try)

		var reason string
		var delay time.Duration
		switch {
		case err != nil:
			if ctx.Err() != nil || !repeatable || retries >= o.MaxRetries {
				return resp, err
			}
			retries++
			reason, delay = "network error: "+err.Error(), o.backoff(retries, nil)
		case resp.StatusCode == http.StatusTooManyRequests ||
			repeatable && (resp.StatusCode == http.StatusBadGateway || resp.StatusCode == http.StatusServiceUnavailable || resp.StatusCode == http.StatusGatewayTimeout):
			if retries >= o.MaxRetries {
				return resp, nil
			}
			retries++
			reason, delay = "status "+strconv.Itoa(resp.StatusCode), o.backoff(retries, resp)
		case repeatable && (resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusConflict):
			if inProgress >= o.InProgressRetries || !hasErrorCode(resp, codeInProgress) {
				return resp, nil
			}
			inProgress++
			delay = defaultInProgressDelay
			if after, ok := o.retryAfter(resp); ok {
				delay = after
			}
			reason, delay = "status "+strconv.Itoa(resp.StatusCode)+" code 12002 (in progress)", min(delay, o.MaxDelay)
		default:
			return resp, nil
		}

		if resp != nil {
			discard(resp)
		}
		if o.OnRetry != nil {
			o.OnRetry(attempt, reason, delay)
		}
		if err := o.sleep(ctx, delay); err != nil {
			return nil, err
		}
	}
}

func (o RetryOptions) withDefaults() RetryOptions {
	if o.MaxRetries == 0 {
		o.MaxRetries = defaultMaxRetries
	}
	if o.BaseDelay <= 0 {
		o.BaseDelay = defaultBaseDelay
	}
	if o.MaxDelay <= 0 {
		o.MaxDelay = defaultMaxDelay
	}
	if o.InProgressRetries == 0 {
		o.InProgressRetries = defaultInProgressRetries
	}
	if o.sleep == nil {
		o.sleep = sleepContext
	}
	if o.random == nil {
		o.random = mathrand.Float64
	}
	if o.now == nil {
		o.now = time.Now
	}
	return o
}

// backoff is the wait before retry n (1-based): full jitter up to min(MaxDelay, BaseDelay * 2^(n-1)), but never
// less than the response's Retry-After, and never more than MaxDelay.
func (o RetryOptions) backoff(n int, resp *http.Response) time.Duration {
	ceiling := o.MaxDelay
	if n-1 < 32 && o.BaseDelay < o.MaxDelay>>(n-1) {
		ceiling = o.BaseDelay << (n - 1)
	}
	delay := time.Duration(o.random() * float64(ceiling))
	if after, ok := o.retryAfter(resp); ok && after > delay {
		delay = after
	}
	return min(delay, o.MaxDelay)
}

// retryAfter reads Retry-After as whole seconds or an HTTP date; an invalid or negative value is ignored.
func (o RetryOptions) retryAfter(resp *http.Response) (time.Duration, bool) {
	if resp == nil {
		return 0, false
	}
	value := strings.TrimSpace(resp.Header.Get("Retry-After"))
	if value == "" {
		return 0, false
	}
	if seconds, err := strconv.ParseInt(value, 10, 64); err == nil {
		switch {
		case seconds < 0:
			return 0, false
		case seconds > int64(o.MaxDelay/time.Second):
			return o.MaxDelay, true // capped anyway; this also keeps a huge value from overflowing
		}
		return time.Duration(seconds) * time.Second, true
	}
	if at, err := http.ParseTime(value); err == nil {
		if wait := at.Sub(o.now()); wait >= 0 {
			return wait, true
		}
	}
	return 0, false
}

// hasErrorCode reads the response body to find its `code`, then puts the body back so the caller can read it.
func hasErrorCode(resp *http.Response, code int) bool {
	if resp.Body == nil || resp.Body == http.NoBody {
		return false
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, maxInspectedBody))
	rest := resp.Body
	resp.Body = struct {
		io.Reader
		io.Closer
	}{io.MultiReader(bytes.NewReader(data), rest), rest}
	if err != nil {
		return false
	}
	var body struct {
		Code json.RawMessage `json:"code"`
	}
	if json.Unmarshal(data, &body) != nil {
		return false
	}
	return strings.Trim(string(body.Code), `"`) == strconv.Itoa(code)
}

// replayableBody returns a function that gives the request body afresh for each attempt (nil when there is no
// body). The generated client builds bodies that already have GetBody; any other body is buffered once.
func replayableBody(req *http.Request) (func() (io.ReadCloser, error), error) {
	if req.Body == nil || req.Body == http.NoBody {
		return nil, nil
	}
	defer req.Body.Close() // a RoundTripper closes the request body; every attempt gets its own copy
	if req.GetBody != nil {
		return req.GetBody, nil
	}
	data, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	return func() (io.ReadCloser, error) { return io.NopCloser(bytes.NewReader(data)), nil }, nil
}

// discard drains and closes a response that is not returned, so its connection can be reused.
func discard(resp *http.Response) {
	if resp.Body != nil {
		_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, maxInspectedBody))
		_ = resp.Body.Close()
	}
}

func sleepContext(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return ctx.Err()
	}
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
