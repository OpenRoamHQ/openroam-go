// Hand-written, not generated (see overlay/). The parts of the OpenRoam client that are not one method per API route:
// options, errors and the helpers the generated resources (openroam.go) call.

package openroam

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Environment is an OpenRoam environment: Production (the default) or Sandbox.
type Environment string

const (
	// Production is https://b2b-api.openroam.com, the default.
	Production Environment = "production"
	// Sandbox is https://b2b-api-sandbox.openroam.com.
	Sandbox Environment = "sandbox"
)

var environments = map[Environment]string{Production: ProductionURL, Sandbox: SandboxURL}

// Option configures New.
type Option func(*clientSettings)

type clientSettings struct {
	environment Environment
	baseURL     string
	maxRetries  int
	httpClient  *http.Client
}

// WithEnvironment picks the environment: Production (the default) or Sandbox.
func WithEnvironment(environment Environment) Option {
	return func(s *clientSettings) { s.environment = environment }
}

// WithBaseURL sends requests to this URL instead of an environment, for example a development host.
func WithBaseURL(baseURL string) Option {
	return func(s *clientSettings) { s.baseURL = baseURL }
}

// WithMaxRetries sets the retries for network errors, 429, 502/503/504 and "still running" (12002). The default is 0: no
// retries. A POST is sent again only with its Idempotency-Key, the same key and the same body. See docs/RETRIES.md.
func WithMaxRetries(maxRetries int) Option {
	return func(s *clientSettings) { s.maxRetries = maxRetries }
}

// WithHTTPClient sends the requests with this client instead of http.DefaultClient. With WithMaxRetries, a copy of it
// whose transport retries is used; the client passed is not changed.
func WithHTTPClient(client *http.Client) Option {
	return func(s *clientSettings) { s.httpClient = client }
}

// New returns a client for the OpenRoam Partner API. An empty apiKey is read from the OPENROAM_API_KEY environment
// variable; when there is none either, New returns an error.
//
//	client, err := openroam.New("your-key", openroam.WithEnvironment(openroam.Sandbox))
func New(apiKey string, opts ...Option) (*OpenRoam, error) {
	settings := clientSettings{environment: Production}
	for _, opt := range opts {
		if opt != nil {
			opt(&settings)
		}
	}
	if apiKey == "" {
		apiKey = os.Getenv("OPENROAM_API_KEY")
	}
	if apiKey == "" {
		return nil, &OpenRoamError{Message: `Missing API key. Pass it as openroam.New("your-key") or set OPENROAM_API_KEY.`}
	}
	baseURL, known := environments[settings.environment]
	if !known {
		return nil, &OpenRoamError{Message: fmt.Sprintf("Unknown environment %q. Use openroam.Production or openroam.Sandbox, or WithBaseURL.", settings.environment)}
	}
	if settings.baseURL != "" {
		baseURL = settings.baseURL
	}
	if settings.maxRetries < 0 {
		return nil, &OpenRoamError{Message: "WithMaxRetries needs 0 or more."}
	}

	httpClient := settings.httpClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	if settings.maxRetries > 0 {
		httpClient = NewRetryingHTTPClient(httpClient, RetryOptions{MaxRetries: settings.maxRetries})
	}
	baseURL = withoutTrailingSlashes(baseURL)
	origin, err := url.Parse(baseURL)
	if err != nil || origin.Scheme == "" || origin.Host == "" {
		return nil, &OpenRoamError{Message: fmt.Sprintf("The base URL %q is not a full URL such as https://host.", baseURL)}
	}
	cfg := NewConfiguration()
	cfg.Servers = ServerConfigurations{{URL: baseURL}}
	cfg.HTTPClient = withPartnerKey(httpClient, apiKey, origin)
	return newOpenRoam(&clientCore{api: NewAPIClient(cfg)}), nil
}

// withPartnerKey returns a copy of client (nil means http.DefaultClient) that sends the key as X-Partner-Key. The key is
// added by the transport, at the last step, so it is never part of the request the generated services build or log. It is
// sent only to the origin the client was made for: a redirect to another host, port or scheme (including https to http)
// goes without it.
func withPartnerKey(client *http.Client, apiKey string, origin *url.URL) *http.Client {
	if client == nil {
		client = http.DefaultClient
	}
	copied := *client
	copied.Transport = &partnerKeyTransport{base: client.Transport, apiKey: apiKey, origin: originOf(origin)}
	return &copied
}

type partnerKeyTransport struct {
	base   http.RoundTripper
	apiKey string
	origin string
}

// originOf is scheme, host and port, with the default port made explicit so "host" and "host:443" are one origin.
func originOf(u *url.URL) string {
	scheme := strings.ToLower(u.Scheme)
	port := u.Port()
	if port == "" {
		port = map[string]string{"http": "80", "https": "443"}[scheme]
	}
	return scheme + "://" + strings.ToLower(u.Hostname()) + ":" + port
}

func (t *partnerKeyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	base := t.base
	if base == nil {
		base = http.DefaultTransport
	}
	keyed := req.Clone(req.Context()) // a RoundTripper must not change the request it is given
	if originOf(req.URL) == t.origin {
		keyed.Header.Set("X-Partner-Key", t.apiKey)
	} else {
		keyed.Header.Del("X-Partner-Key")
	}
	return base.RoundTrip(keyed)
}

// CallOption configures one call.
type CallOption func(*callSettings)

type callSettings struct {
	key    string
	hasKey bool
}

// WithIdempotencyKey sets the Idempotency-Key of a call that creates something: an order, a refund and the like. If the
// same request arrives twice with the same key, OpenRoam does it once. Without it, every such call gets a new random key,
// which makes a retry by this SDK safe. To stay safe even when your own program crashes and asks again, pass your own,
// made once per purchase and kept, for example the id of your order. Calls that need no key ignore it.
func WithIdempotencyKey(key string) CallOption {
	return func(s *callSettings) { s.key, s.hasKey = key, true }
}

func (s callSettings) idempotencyKey() string {
	if s.hasKey {
		return s.key
	}
	return NewIdempotencyKey()
}

// clientCore is what every resource of one client shares.
type clientCore struct {
	api *APIClient
}

// begin reads the call options.
func (c *clientCore) begin(ctx context.Context, opts []CallOption) (context.Context, callSettings, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var settings callSettings
	for _, opt := range opts {
		if opt != nil {
			opt(&settings)
		}
	}
	if settings.hasKey && settings.key == "" {
		return ctx, settings, &OpenRoamError{Message: "WithIdempotencyKey needs a key; it is empty."}
	}
	return ctx, settings, nil
}

// OpenRoamError is the error of every failed call. Use errors.As to read it.
//
// For an error answer from the API it carries the error body as the API sends it: Status (the HTTP status), Code (the
// number to branch on), Name (its machine name, for example INVALID_REQUEST), Message, RequestID (keep it for support)
// and Details (per-field problems). When no answer arrived (network error, timeout, a cancelled context), Status is 0 and
// Unwrap returns the original error, so errors.Is(err, context.Canceled) and errors.Is(err, context.DeadlineExceeded)
// still work.
type OpenRoamError struct {
	Status    int
	Code      int
	Name      string
	Message   string
	RequestID string
	Details   []ErrorDetail
	cause     error
}

func (e *OpenRoamError) Error() string {
	var b strings.Builder
	b.WriteString("openroam: ")
	if e.Status != 0 {
		fmt.Fprintf(&b, "status %d", e.Status)
		if e.Code != 0 {
			fmt.Fprintf(&b, ", code %d", e.Code)
		}
		if e.Name != "" {
			fmt.Fprintf(&b, " %s", e.Name)
		}
		b.WriteString(": ")
	}
	b.WriteString(e.Message)
	if e.RequestID != "" {
		fmt.Fprintf(&b, " (request %s)", e.RequestID)
	}
	return b.String()
}

// Unwrap returns the error that made the call fail when no answer (or only part of one) arrived, or nil.
func (e *OpenRoamError) Unwrap() error { return e.cause }

// finish turns what a generated Execute returned into the client's result: an error answer, a 202 that says the
// request is still running (12002) and a failure without an answer all become an *OpenRoamError.
func finish[T any](value T, resp *http.Response, err error) (T, error) {
	if failure := failed(resp, err); failure != nil {
		var zero T
		return zero, failure
	}
	return value, nil
}

// finishWithoutValue is finish for a route that answers without a body.
func finishWithoutValue(resp *http.Response, err error) error {
	return failed(resp, err)
}

func failed(resp *http.Response, err error) error {
	if resp != nil && (resp.StatusCode >= 300 || resp.StatusCode == http.StatusAccepted) {
		var body []byte
		if resp.Body != nil {
			body, _ = io.ReadAll(io.LimitReader(resp.Body, maxInspectedBody))
		}
		e := fromBody(resp.StatusCode, body)
		if resp.StatusCode >= 300 || e.Code == codeInProgress {
			// The answer arrived, but a cancelled or expired ctx may have cut its body short: keep that as the cause, so
			// errors.Is(err, ctx.Err()) still holds.
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				e.cause = err
			}
			return e
		}
	}
	if err == nil {
		return nil
	}
	var answerless *url.Error
	if errors.As(err, &answerless) {
		return &OpenRoamError{Message: "The request did not get an answer: " + answerless.Err.Error(), cause: err}
	}
	e := &OpenRoamError{Message: err.Error(), cause: err}
	if resp != nil {
		e.Status = resp.StatusCode
	}
	return e
}

// fromBody reads an error body. Each field is read on its own, so one odd field never hides the others; an empty or
// non-JSON body (a proxy error page, say) reports the status only.
func fromBody(status int, body []byte) *OpenRoamError {
	e := &OpenRoamError{Status: status, Message: fmt.Sprintf("The request failed with status %d.", status)}
	var fields map[string]json.RawMessage
	if json.Unmarshal(body, &fields) != nil {
		return e
	}
	read := func(name string, into any) bool {
		raw := fields[name]
		return raw != nil && strings.TrimSpace(string(raw)) != "null" && json.Unmarshal(raw, into) == nil
	}
	var code int
	if read("code", &code) {
		e.Code = code
	}
	var text string
	if read("name", &text) {
		e.Name = text
	}
	if read("message", &text) {
		e.Message = text
	}
	if read("requestId", &text) {
		e.RequestID = text
	}
	var details []ErrorDetail
	if read("details", &details) {
		e.Details = details
	}
	return e
}

// A loop, not a regular expression: the URL is caller input.
func withoutTrailingSlashes(s string) string {
	end := len(s)
	for end > 0 && s[end-1] == '/' {
		end--
	}
	return s[:end]
}
