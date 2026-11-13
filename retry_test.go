package openroam

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"
)

// One scripted answer of the test server.
type answer struct {
	status     int
	retryAfter string
	body       string
}

// A request as the test server received it.
type received struct {
	method, key, body string
}

type retryServer struct {
	*httptest.Server
	mu       sync.Mutex
	answers  []answer
	requests []received
}

func newRetryServer(t *testing.T, answers ...answer) *retryServer {
	s := &retryServer{answers: answers}
	s.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		s.mu.Lock()
		defer s.mu.Unlock()
		s.requests = append(s.requests, received{r.Method, r.Header.Get("Idempotency-Key"), string(body)})
		if len(s.requests) > len(s.answers) {
			t.Errorf("unexpected request %d", len(s.requests))
			w.WriteHeader(http.StatusTeapot)
			return
		}
		a := s.answers[len(s.requests)-1]
		if a.retryAfter != "" {
			w.Header().Set("Retry-After", a.retryAfter)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(a.status)
		_, _ = io.WriteString(w, a.body)
	}))
	t.Cleanup(s.Close)
	return s
}

func (s *retryServer) seen() []received {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]received(nil), s.requests...)
}

// testOptions waits for no real time and records every wait.
func testOptions(opts RetryOptions, waits *[]time.Duration) RetryOptions {
	opts.sleep = func(ctx context.Context, d time.Duration) error {
		*waits = append(*waits, d)
		return ctx.Err()
	}
	if opts.random == nil {
		opts.random = func() float64 { return 0.5 }
	}
	return opts
}

func send(t *testing.T, client *http.Client, method, url, key, body string) (*http.Response, string) {
	t.Helper()
	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		t.Fatal(err)
	}
	if key != "" {
		req.Header.Set("Idempotency-Key", key)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	return resp, string(data)
}

const (
	inProgress = `{"code":12002,"name":"IDEMPOTENCY_KEY_IN_PROGRESS","message":"still running"}`
	keyReused  = `{"code":12001,"name":"IDEMPOTENCY_KEY_REUSED","message":"different body"}`
	orderBody  = `{"items":[{"publicId":"offer_1","inputs":{"phoneNumber":"+15550000000"}}]}`
)

func expectRequests(t *testing.T, s *retryServer, n int) []received {
	t.Helper()
	got := s.seen()
	if len(got) != n {
		t.Fatalf("got %d requests, want %d: %+v", len(got), n, got)
	}
	return got
}

func expectSameRequest(t *testing.T, got []received) {
	t.Helper()
	for _, r := range got[1:] {
		if r != got[0] {
			t.Errorf("a retry differs from the first request: %+v vs %+v", r, got[0])
		}
	}
	if got[0].key == "" || got[0].body == "" {
		t.Errorf("the request lost its key or body: %+v", got[0])
	}
}

// 1
func TestRetryGet503Then200(t *testing.T) {
	s := newRetryServer(t, answer{status: 503}, answer{status: 200, body: `{"ok":true}`})
	var waits []time.Duration
	resp, body := send(t, NewRetryingHTTPClient(nil, testOptions(RetryOptions{}, &waits)), "GET", s.URL, "", "")
	if resp.StatusCode != 200 || body != `{"ok":true}` {
		t.Fatalf("got %d %s", resp.StatusCode, body)
	}
	expectRequests(t, s, 2)
	if len(waits) != 1 || waits[0] > 500*time.Millisecond {
		t.Errorf("waits %v, want one within the first backoff step", waits)
	}
}

// 2
func TestRetryGet429HonoursRetryAfterCappedAtMaxDelay(t *testing.T) {
	s := newRetryServer(t, answer{status: 429, retryAfter: "3"}, answer{status: 429, retryAfter: "120"}, answer{status: 200})
	var waits []time.Duration
	resp, _ := send(t, NewRetryingHTTPClient(nil, testOptions(RetryOptions{MaxDelay: 10 * time.Second}, &waits)), "GET", s.URL, "", "")
	if resp.StatusCode != 200 {
		t.Fatalf("got %d", resp.StatusCode)
	}
	expectRequests(t, s, 3)
	if len(waits) != 2 || waits[0] != 3*time.Second || waits[1] != 10*time.Second {
		t.Errorf("waits %v, want [3s 10s]", waits)
	}
}

// 3
func TestRetryPostWithKey503Then200SendsTheSameKeyAndBody(t *testing.T) {
	s := newRetryServer(t, answer{status: 503}, answer{status: 200})
	var waits []time.Duration
	resp, _ := send(t, NewRetryingHTTPClient(nil, testOptions(RetryOptions{}, &waits)), "POST", s.URL, "key-1", orderBody)
	if resp.StatusCode != 200 {
		t.Fatalf("got %d", resp.StatusCode)
	}
	expectSameRequest(t, expectRequests(t, s, 2))
}

// 4 and 5
func TestRetryPostWithKeyWhileInProgress(t *testing.T) {
	for _, status := range []int{202, 409} {
		s := newRetryServer(t, answer{status: status, body: inProgress, retryAfter: "1"}, answer{status: status, body: inProgress}, answer{status: 200})
		var waits []time.Duration
		var reasons []string
		opts := testOptions(RetryOptions{OnRetry: func(attempt int, reason string, delay time.Duration) {
			reasons = append(reasons, reason)
		}}, &waits)
		resp, _ := send(t, NewRetryingHTTPClient(nil, opts), "POST", s.URL, "key-1", orderBody)
		if resp.StatusCode != 200 {
			t.Fatalf("%d: got %d", status, resp.StatusCode)
		}
		expectSameRequest(t, expectRequests(t, s, 3))
		if len(waits) != 2 || waits[0] != time.Second || waits[1] != 2*time.Second {
			t.Errorf("%d: waits %v, want [1s 2s] (Retry-After, then the 2 s default)", status, waits)
		}
		if len(reasons) != 2 || !strings.Contains(reasons[0], "12002") {
			t.Errorf("%d: OnRetry reasons %v", status, reasons)
		}
	}
}

func TestRetryInProgressHasItsOwnBudget(t *testing.T) {
	answers := []answer{{status: 503}, {status: 503}, {status: 503}}
	for i := 0; i < 4; i++ {
		answers = append(answers, answer{status: 202, body: inProgress})
	}
	s := newRetryServer(t, append(answers, answer{status: 200})...)
	var waits []time.Duration
	resp, _ := send(t, NewRetryingHTTPClient(nil, testOptions(RetryOptions{InProgressRetries: 4}, &waits)), "POST", s.URL, "key-1", orderBody)
	if resp.StatusCode != 200 {
		t.Fatalf("got %d", resp.StatusCode)
	}
	expectRequests(t, s, 8)

	s = newRetryServer(t, answer{status: 202, body: inProgress}, answer{status: 202, body: inProgress})
	resp, body := send(t, NewRetryingHTTPClient(nil, testOptions(RetryOptions{InProgressRetries: 1}, &waits)), "POST", s.URL, "key-1", orderBody)
	if resp.StatusCode != 202 || body != inProgress {
		t.Fatalf("the last in-progress answer is returned unchanged; got %d %s", resp.StatusCode, body)
	}
	expectRequests(t, s, 2)
}

// 6, 7 and 9: never repeated.
func TestRetryNeverRepeats(t *testing.T) {
	cases := []struct {
		name   string
		method string
		key    string
		answer answer
	}{
		{"409 12001", "POST", "key-1", answer{status: 409, body: keyReused}},
		{"409 12000", "POST", "key-1", answer{status: 409, body: `{"code":12000}`}},
		{"500", "POST", "key-1", answer{status: 500, body: `{"code":10000}`}},
		{"202 accepted", "POST", "key-1", answer{status: 202, body: `{"id":"order_1"}`}},
		{"400", "GET", "", answer{status: 400, body: `{"code":10001}`}},
		{"401", "GET", "", answer{status: 401}},
		{"404", "POST", "key-1", answer{status: 404}},
	}
	for _, c := range cases {
		s := newRetryServer(t, c.answer)
		var waits []time.Duration
		resp, body := send(t, NewRetryingHTTPClient(nil, testOptions(RetryOptions{}, &waits)), c.method, s.URL, c.key, orderBody)
		if resp.StatusCode != c.answer.status || body != c.answer.body {
			t.Errorf("%s: got %d %s", c.name, resp.StatusCode, body)
		}
		expectRequests(t, s, 1)
		if len(waits) != 0 {
			t.Errorf("%s: waited %v", c.name, waits)
		}
	}
}

// 8
func TestRetryPostWithoutKey(t *testing.T) {
	s := newRetryServer(t, answer{status: 503, body: `{"code":1}`})
	var waits []time.Duration
	client := NewRetryingHTTPClient(nil, testOptions(RetryOptions{}, &waits))
	resp, body := send(t, client, "POST", s.URL, "", orderBody)
	if resp.StatusCode != 503 || body != `{"code":1}` {
		t.Fatalf("got %d %s", resp.StatusCode, body)
	}
	expectRequests(t, s, 1)

	s = newRetryServer(t, answer{status: 429, retryAfter: "2"}, answer{status: 200})
	resp, _ = send(t, client, "POST", s.URL, "", orderBody)
	if resp.StatusCode != 200 {
		t.Fatalf("a refused request is retried; got %d", resp.StatusCode)
	}
	got := expectRequests(t, s, 2)
	if got[0].body != orderBody || got[1].body != orderBody {
		t.Errorf("the body changed: %+v", got)
	}

	// A 202 with 12002 is not looked at for a request that is not repeatable.
	s = newRetryServer(t, answer{status: 202, body: inProgress})
	if resp, _ = send(t, client, "POST", s.URL, "", orderBody); resp.StatusCode != 202 {
		t.Fatalf("got %d", resp.StatusCode)
	}
	expectRequests(t, s, 1)
}

// flaky fails its first n round trips with a network error.
type flaky struct {
	n     int
	calls int
}

func (f *flaky) RoundTrip(req *http.Request) (*http.Response, error) {
	f.calls++
	if f.calls <= f.n {
		if req.Body != nil {
			req.Body.Close()
		}
		return nil, errors.New("connection reset by peer")
	}
	return http.DefaultTransport.RoundTrip(req)
}

// 10
func TestRetryNetworkErrorThenSuccess(t *testing.T) {
	s := newRetryServer(t, answer{status: 200, body: `{"ok":true}`}, answer{status: 200})
	var waits []time.Duration
	base := &flaky{n: 1}
	client := &http.Client{Transport: &RetryTransport{Base: base, Options: testOptions(RetryOptions{}, &waits)}}
	resp, body := send(t, client, "GET", s.URL, "", "")
	if resp.StatusCode != 200 || body != `{"ok":true}` || base.calls != 2 {
		t.Fatalf("got %d %s after %d calls", resp.StatusCode, body, base.calls)
	}

	// The same error is surfaced at once for a POST without a key.
	base = &flaky{n: 1}
	client = &http.Client{Transport: &RetryTransport{Base: base, Options: testOptions(RetryOptions{}, &waits)}}
	req, _ := http.NewRequest("POST", s.URL, strings.NewReader(orderBody))
	if _, err := client.Do(req); err == nil || base.calls != 1 {
		t.Fatalf("want the network error after one call; got %v after %d", err, base.calls)
	}
}

// 11
func TestRetryExhaustedReturnsTheLastResponse(t *testing.T) {
	answers := []answer{{status: 503, body: "1"}, {status: 502, body: "2"}, {status: 504, body: "3"}, {status: 503, body: "last"}}
	s := newRetryServer(t, answers...)
	var waits []time.Duration
	var attempts []int
	opts := testOptions(RetryOptions{BaseDelay: 100 * time.Millisecond, OnRetry: func(attempt int, _ string, _ time.Duration) {
		attempts = append(attempts, attempt)
	}}, &waits)
	opts.random = func() float64 { return 0.999999 }
	resp, body := send(t, NewRetryingHTTPClient(nil, opts), "POST", s.URL, "key-1", orderBody)
	if resp.StatusCode != 503 || body != "last" {
		t.Fatalf("got %d %s", resp.StatusCode, body)
	}
	expectSameRequest(t, expectRequests(t, s, 3+1))
	// Full jitter below 100, 200 and 400 ms.
	for i, ceiling := range []time.Duration{100, 200, 400} {
		if waits[i] > ceiling*time.Millisecond || waits[i] < ceiling*time.Millisecond*9/10 {
			t.Errorf("wait %d is %v, want just under %v", i+1, waits[i], ceiling*time.Millisecond)
		}
	}
	if len(attempts) != 3 || attempts[0] != 1 || attempts[2] != 3 {
		t.Errorf("OnRetry attempts %v", attempts)
	}
}

// 12
func TestRetryAfterParsing(t *testing.T) {
	now := time.Date(2026, 11, 4, 12, 0, 0, 0, time.UTC)
	o := RetryOptions{now: func() time.Time { return now }}.withDefaults()
	for value, want := range map[string]time.Duration{
		"5":    5 * time.Second,
		"0":    0,
		"3600": 30 * time.Second, // capped at MaxDelay
		now.Add(7 * time.Second).Format(http.TimeFormat): 7 * time.Second,
	} {
		got, ok := o.retryAfter(&http.Response{Header: http.Header{"Retry-After": {value}}})
		if !ok || got != want {
			t.Errorf("%q: got %v %v, want %v", value, got, ok, want)
		}
	}
	for _, value := range []string{"-1", "soon", "1.5", now.Add(-time.Minute).Format(http.TimeFormat)} {
		if got, ok := o.retryAfter(&http.Response{Header: http.Header{"Retry-After": {value}}}); ok {
			t.Errorf("%q: got %v, want it ignored", value, got)
		}
	}

	// An invalid value falls back to the backoff; a date is honoured end to end.
	s := newRetryServer(t, answer{status: 503, retryAfter: "soon"}, answer{status: 503, retryAfter: now.Add(4 * time.Second).Format(http.TimeFormat)}, answer{status: 200})
	var waits []time.Duration
	opts := testOptions(RetryOptions{}, &waits)
	opts.now = func() time.Time { return now }
	if resp, _ := send(t, NewRetryingHTTPClient(nil, opts), "GET", s.URL, "", ""); resp.StatusCode != 200 {
		t.Fatalf("got %d", resp.StatusCode)
	}
	if len(waits) != 2 || waits[0] != 250*time.Millisecond || waits[1] != 4*time.Second {
		t.Errorf("waits %v, want [250ms 4s]", waits)
	}
}

// 13
func TestRetryNewIdempotencyKey(t *testing.T) {
	uuid4 := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	seen := map[string]bool{}
	for i := 0; i < 1000; i++ {
		key := NewIdempotencyKey()
		if !uuid4.MatchString(key) || seen[key] {
			t.Fatalf("not a new UUID v4: %s", key)
		}
		seen[key] = true
	}
}

func TestRetryStopsWhenTheContextIsCancelled(t *testing.T) {
	s := newRetryServer(t, answer{status: 503}, answer{status: 200})
	ctx, cancel := context.WithCancel(context.Background())
	client := NewRetryingHTTPClient(nil, RetryOptions{BaseDelay: time.Minute, MaxDelay: time.Hour, OnRetry: func(int, string, time.Duration) { cancel() }})
	req, _ := http.NewRequestWithContext(ctx, "GET", s.URL, nil)
	start := time.Now()
	if _, err := client.Do(req); !errors.Is(err, context.Canceled) {
		t.Fatalf("got %v", err)
	}
	if time.Since(start) > 5*time.Second {
		t.Error("the wait was not cut short")
	}
	expectRequests(t, s, 1)
}

func TestRetryKeepsTheBaseClientSettings(t *testing.T) {
	base := &http.Client{Timeout: 7 * time.Second}
	client := NewRetryingHTTPClient(base, RetryOptions{})
	if client == base || client.Timeout != 7*time.Second || base.Transport != nil {
		t.Fatal("the base client is copied, not changed")
	}
}

// The generated client builds a request that can be sent again byte for byte.
func TestRetryThroughTheGeneratedClient(t *testing.T) {
	s := newRetryServer(t, answer{status: 503}, answer{status: 202, body: inProgress}, answer{status: 200, body: `{}`})
	var waits []time.Duration
	cfg := NewConfiguration()
	cfg.Servers = ServerConfigurations{{URL: s.URL}}
	cfg.HTTPClient = NewRetryingHTTPClient(nil, testOptions(RetryOptions{}, &waits))
	client := NewAPIClient(cfg)
	ctx := context.WithValue(context.Background(), ContextAPIKeys, map[string]APIKey{"PartnerApiKey": {Key: "key"}})

	item := NewAirtimeOrderRequestItem("offer_1", *NewAirtimeCartInputs("+15550000000"))
	_, httpResp, _ := client.AirtimeAPI.AirtimeCreateOrder(ctx).
		IdempotencyKey(NewIdempotencyKey()).
		AirtimeOrderRequest(*NewAirtimeOrderRequest([]AirtimeOrderRequestItem{*item})).
		Execute()
	if httpResp == nil || httpResp.StatusCode != 200 {
		t.Fatalf("got %+v", httpResp)
	}
	got := expectRequests(t, s, 3)
	expectSameRequest(t, got)
	if !strings.Contains(got[0].body, "offer_1") {
		t.Errorf("body %s", got[0].body)
	}
}
