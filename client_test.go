package openroam

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"
)

var uuidV4 = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

// One scripted answer of the client test server; drop closes the connection without an answer.
type reply struct {
	status  int
	headers map[string]string
	body    string
	drop    bool
}

type clientRequest struct {
	method, url, body string
	header            http.Header
}

// scriptedServer answers each request with the next scripted reply (the last one repeats) and records what it got.
func scriptedServer(t *testing.T, script ...reply) (string, func() []clientRequest) {
	t.Helper()
	var mu sync.Mutex
	var requests []clientRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		mu.Lock()
		requests = append(requests, clientRequest{r.Method, r.URL.RequestURI(), string(body), r.Header.Clone()})
		next := script[min(len(requests), len(script))-1]
		mu.Unlock()
		if next.drop {
			conn, _, err := w.(http.Hijacker).Hijack()
			if err == nil {
				_ = conn.Close()
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		for k, v := range next.headers {
			w.Header().Set(k, v)
		}
		w.WriteHeader(next.status)
		body = []byte(next.body)
		if next.body == "" {
			body = []byte("{}")
		}
		_, _ = w.Write(body)
	}))
	t.Cleanup(server.Close)
	return server.URL, func() []clientRequest {
		mu.Lock()
		defer mu.Unlock()
		return slices.Clone(requests)
	}
}

func mustNew(t *testing.T, apiKey string, opts ...Option) *OpenRoam {
	t.Helper()
	c, err := New(apiKey, opts...)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func asOpenRoamError(t *testing.T, err error) *OpenRoamError {
	t.Helper()
	var e *OpenRoamError
	if !errors.As(err, &e) {
		t.Fatalf("want an *OpenRoamError, got %T: %v", err, err)
	}
	return e
}

var (
	okReply       = reply{status: 200, body: `{"items":[],"pagination":{}}`}
	esimOrderJSON = `{"currency":"USD","items":[{"publicId":"esim_1","quantity":1,"inputs":{}}]}`
)

func esimOrder(t *testing.T) *EsimOrderRequest {
	t.Helper()
	var body EsimOrderRequest
	readSample(t, []byte(esimOrderJSON), &body)
	return &body
}

func sameJSON(t *testing.T, got, want string) {
	t.Helper()
	var a, b any
	if err := json.Unmarshal([]byte(got), &a); err != nil {
		t.Fatalf("not JSON: %q", got)
	}
	_ = json.Unmarshal([]byte(want), &b)
	if !reflect.DeepEqual(a, b) {
		t.Errorf("body %s, want %s", got, want)
	}
}

func TestASearchIsOneCallOnTheProductWithTheKeyAsXPartnerKey(t *testing.T) {
	baseURL, requests := scriptedServer(t, okReply)
	found, err := mustNew(t, "secret", WithBaseURL(baseURL)).Esims.Search(context.Background(), &EsimsSearchQuery{CountryCode: PtrString("JP")})
	if err != nil {
		t.Fatal(err)
	}
	if found == nil || len(found.Items) != 0 {
		t.Errorf("found %+v", found)
	}
	r := requests()[0]
	if r.method != "GET" || r.url != "/api/v1/b2b/esims/search?countryCode=JP" || r.header.Get("X-Partner-Key") != "secret" {
		t.Errorf("request %s %s, key %q", r.method, r.url, r.header.Get("X-Partner-Key"))
	}
}

func TestTheKeyIsAddedByTheTransportSoTheGeneratedDebugLogNeverHoldsIt(t *testing.T) {
	baseURL, requests := scriptedServer(t, okReply)
	transport := &recordingTransport{}
	own := &http.Client{Transport: transport}
	c := mustNew(t, "secret-key", WithBaseURL(baseURL))
	c.Esims.c.api.GetConfig().Debug = true
	var logged bytes.Buffer
	log.SetOutput(&logged)
	defer log.SetOutput(os.Stderr)
	if _, err := c.Esims.Search(context.Background(), nil); err != nil {
		t.Fatal(err)
	}
	if got := requests()[0].header.Get("X-Partner-Key"); got != "secret-key" {
		t.Errorf("X-Partner-Key %q", got)
	}
	if !strings.Contains(logged.String(), "/api/v1/b2b/esims/search") || strings.Contains(logged.String(), "secret-key") {
		t.Errorf("debug log %q", logged.String())
	}
	// Your own client is used through a copy and stays as it was.
	if _, err := mustNew(t, "k", WithHTTPClient(own), WithMaxRetries(1)).Esims.Search(context.Background(), nil); err != nil {
		t.Fatal(err)
	}
	if own.Transport != transport || len(transport.urls) != 1 {
		t.Errorf("your client was changed or not used: %T, %v", own.Transport, transport.urls)
	}
}

func TestTheKeyGoesOnlyToTheConfiguredOriginEvenAcrossARedirect(t *testing.T) {
	var elsewhereHeaders []http.Header
	elsewhere := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		elsewhereHeaders = append(elsewhereHeaders, r.Header.Clone())
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"items":[],"pagination":{}}`))
	}))
	defer elsewhere.Close()
	var configuredHeaders []http.Header
	configured := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		configuredHeaders = append(configuredHeaders, r.Header.Clone())
		switch r.URL.Path {
		case "/api/v1/b2b/esims/search":
			http.Redirect(w, r, elsewhere.URL+"/stolen", http.StatusFound)
		default:
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"items":[],"pagination":{}}`))
		}
	}))
	defer configured.Close()

	client := mustNew(t, "secret", WithBaseURL(configured.URL))
	if _, err := client.Esims.Search(context.Background(), nil); err != nil {
		t.Fatalf("search: %v", err)
	}

	if len(configuredHeaders) != 1 || configuredHeaders[0].Get("X-Partner-Key") != "secret" {
		t.Fatalf("the configured origin must get the key, got %v", configuredHeaders)
	}
	if len(elsewhereHeaders) != 1 {
		t.Fatalf("the redirect must be followed once, got %d requests", len(elsewhereHeaders))
	}
	if got := elsewhereHeaders[0].Get("X-Partner-Key"); got != "" {
		t.Fatalf("the key must not follow a redirect to another origin, got %q", got)
	}
}

func TestOneOriginMeansTheSameSchemeHostAndPortWithDefaultPortsMadeExplicit(t *testing.T) {
	cases := []struct {
		a, b string
		same bool
	}{
		{"https://b2b-api.openroam.com", "https://b2b-api.openroam.com:443", true},
		{"https://B2B-API.openroam.com/x", "https://b2b-api.openroam.com/y", true},
		{"https://b2b-api.openroam.com", "http://b2b-api.openroam.com", false},
		{"https://b2b-api.openroam.com", "https://b2b-api-sandbox.openroam.com", false},
		{"http://127.0.0.1:8080", "http://127.0.0.1:8081", false},
	}
	for _, c := range cases {
		a, _ := url.Parse(c.a)
		b, _ := url.Parse(c.b)
		if got := originOf(a) == originOf(b); got != c.same {
			t.Errorf("%s vs %s: same origin = %v, want %v", c.a, c.b, got, c.same)
		}
	}
}

func TestABaseURLThatIsNotAFullURLIsRefused(t *testing.T) {
	if _, err := New("key", WithBaseURL("openroam.test")); err == nil {
		t.Fatal("expected an error for a base URL without a scheme")
	}
}

func TestTheKeyFallsBackToOPENROAM_API_KEY(t *testing.T) {
	if c, err := New("key"); err != nil || c.Esims == nil {
		t.Fatalf("New(key): %v", err)
	}
	t.Setenv("OPENROAM_API_KEY", "from-env")
	baseURL, requests := scriptedServer(t, okReply)
	if _, err := mustNew(t, "", WithBaseURL(baseURL)).Esims.Search(context.Background(), nil); err != nil {
		t.Fatal(err)
	}
	if got := requests()[0].header.Get("X-Partner-Key"); got != "from-env" {
		t.Errorf("X-Partner-Key %q", got)
	}
	t.Setenv("OPENROAM_API_KEY", "")
	_, err := New("")
	if e := asOpenRoamError(t, err); !strings.Contains(e.Message, "Missing API key") {
		t.Errorf("message %q", e.Message)
	}
}

type recordingTransport struct{ urls []string }

func (r *recordingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	r.urls = append(r.urls, req.URL.String())
	return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": {"application/json"}},
		Body: io.NopCloser(strings.NewReader("{}")), Request: req}, nil
}

func TestTheEnvironmentPicksTheURLAndAnUnknownOneIsRefused(t *testing.T) {
	transport := &recordingTransport{}
	client := &http.Client{Transport: transport}
	ctx := context.Background()
	if _, err := mustNew(t, "k", WithHTTPClient(client)).Esims.Search(ctx, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := mustNew(t, "k", WithHTTPClient(client), WithEnvironment(Sandbox)).Esims.Search(ctx, nil); err != nil {
		t.Fatal(err)
	}
	if transport.urls[0] != "https://b2b-api.openroam.com/api/v1/b2b/esims/search" || !strings.HasPrefix(transport.urls[1], "https://b2b-api-sandbox.openroam.com/") {
		t.Errorf("urls %v", transport.urls)
	}
	_, err := New("k", WithEnvironment("staging"))
	if e := asOpenRoamError(t, err); !strings.Contains(e.Message, "Unknown environment") {
		t.Errorf("message %q", e.Message)
	}
	if _, err := New("k", WithMaxRetries(-1)); err == nil {
		t.Error("a negative WithMaxRetries is accepted")
	}
}

func TestATrailingSlashOnTheBaseURLIsIgnored(t *testing.T) {
	baseURL, requests := scriptedServer(t, okReply)
	if _, err := mustNew(t, "k", WithBaseURL(baseURL+"///")).Esims.Search(context.Background(), nil); err != nil {
		t.Fatal(err)
	}
	if got := requests()[0].url; got != "/api/v1/b2b/esims/search" {
		t.Errorf("url %s", got)
	}
}

func TestAQuoteSendsTheBodyAsGivenNotWrapped(t *testing.T) {
	baseURL, requests := scriptedServer(t, reply{status: 200, body: `{"currency":"USD"}`})
	if _, err := mustNew(t, "k", WithBaseURL(baseURL)).Esims.Quote(context.Background(), esimOrder(t)); err != nil {
		t.Fatal(err)
	}
	r := requests()[0]
	if r.method != "POST" || r.url != "/api/v1/b2b/esims/quote" {
		t.Errorf("request %s %s", r.method, r.url)
	}
	sameJSON(t, r.body, esimOrderJSON)
	if r.header.Get("Idempotency-Key") != "" {
		t.Errorf("a quote sent an Idempotency-Key")
	}
}

func TestAnOrderGetsItsOwnIdempotencyKeyAndYourOwnKeyWins(t *testing.T) {
	baseURL, requests := scriptedServer(t, okReply)
	c := mustNew(t, "k", WithBaseURL(baseURL))
	ctx := context.Background()
	for _, opts := range [][]CallOption{nil, nil, {WithIdempotencyKey("my-order-42")}} {
		if _, err := c.Esims.Order(ctx, esimOrder(t), opts...); err != nil {
			t.Fatal(err)
		}
	}
	got := requests()
	first, second, third := got[0].header.Get("Idempotency-Key"), got[1].header.Get("Idempotency-Key"), got[2].header.Get("Idempotency-Key")
	if !uuidV4.MatchString(first) || first == second || third != "my-order-42" {
		t.Errorf("keys %q %q %q", first, second, third)
	}
	sameJSON(t, got[0].body, esimOrderJSON)

	_, err := c.Esims.Order(ctx, esimOrder(t), WithIdempotencyKey(""))
	if e := asOpenRoamError(t, err); e.Status != 0 || len(requests()) != 3 {
		t.Errorf("an empty key: %v, %d requests", e, len(requests()))
	}
}

func TestPathIdsComeFirstInPathOrderAndQueryParametersAreOneStruct(t *testing.T) {
	baseURL, requests := scriptedServer(t, reply{status: 200})
	c := mustNew(t, "k", WithBaseURL(baseURL))
	ctx := context.Background()
	_, _ = c.Esims.GetUsage(ctx, "8988-1")
	_, _ = c.OrderItems.GetAmendment(ctx, "item-1", "amend-2")
	_, _ = c.Esims.ListUsageHistory(ctx, "8988-1", &EsimsListUsageHistoryQuery{Limit: PtrInt32(5)})
	// A query parameter left out that has a documented default (page) is sent with that default by the generated service.
	want := []string{"/api/v1/b2b/esims/8988-1/usage", "/api/v1/b2b/order-items/item-1/amendments/amend-2", "/api/v1/b2b/esims/8988-1/usage/history?limit=5&page=1"}
	for i, r := range requests() {
		if r.url != want[i] {
			t.Errorf("url %s, want %s", r.url, want[i])
		}
	}
}

func TestAQueryValueCannotReplaceAPositionalId(t *testing.T) {
	// The query struct has no field for a path id, so it cannot carry one.
	for _, query := range []any{EsimsGetDestinationQuery{}, EsimsListUsageHistoryQuery{}, AirtimeGetOperatorQuery{}, WalletsListTransactionsQuery{}} {
		for _, id := range []string{"Id", "ID", "Iccid", "WalletId"} {
			if _, has := reflect.TypeOf(query).FieldByName(id); has {
				t.Errorf("%T has a field %s", query, id)
			}
		}
	}
	baseURL, requests := scriptedServer(t, reply{status: 200})
	_, _ = mustNew(t, "k", WithBaseURL(baseURL)).Esims.GetDestination(context.Background(), "A", &EsimsGetDestinationQuery{Currency: PtrString("USD")})
	if got := requests()[0].url; got != "/api/v1/b2b/esims/destinations/A?currency=USD" {
		t.Errorf("url %s", got)
	}
}

func TestTheNestedOrderItemsManualAmendmentsResourcePostsWithAKey(t *testing.T) {
	baseURL, requests := scriptedServer(t, okReply)
	_, _ = mustNew(t, "k", WithBaseURL(baseURL)).OrderItems.ManualAmendments.Create(context.Background(), "item-1",
		&CreateB2BManualAmendmentCommand{Attributes: map[string]any{"checkInDate": "2026-01-02"}})
	r := requests()[0]
	if r.method != "POST" || r.url != "/api/v1/b2b/order-items/item-1/manual-amendments" || r.header.Get("Idempotency-Key") == "" {
		t.Errorf("request %s %s, key %q", r.method, r.url, r.header.Get("Idempotency-Key"))
	}
	sameJSON(t, r.body, `{"attributes":{"checkInDate":"2026-01-02"}}`)
}

func TestAFailedCallReturnsAnOpenRoamErrorCarryingTheErrorBody(t *testing.T) {
	baseURL, _ := scriptedServer(t, reply{status: 409, body: `{"code":11002,"name":"PRICE_ABOVE_LIMIT","message":"The price is above your limit.",` +
		`"requestId":"req_9","details":[{"path":"items[0].inputs.maxPrice","reason":"PRICE_ABOVE_LIMIT"}]}`})
	result, err := mustNew(t, "k", WithBaseURL(baseURL)).Esims.Quote(context.Background(), esimOrder(t))
	e := asOpenRoamError(t, err)
	if result != nil || e.Status != 409 || e.Code != 11002 || e.Name != "PRICE_ABOVE_LIMIT" || e.Message != "The price is above your limit." || e.RequestID != "req_9" {
		t.Errorf("error %+v, result %v", e, result)
	}
	if len(e.Details) != 1 || e.Details[0].GetPath() != "items[0].inputs.maxPrice" || e.Details[0].GetReason() != "PRICE_ABOVE_LIMIT" {
		t.Errorf("details %+v", e.Details)
	}
	if e.Unwrap() != nil {
		t.Errorf("an error answer has a cause: %v", e.Unwrap())
	}
	if !strings.Contains(err.Error(), "409") || !strings.Contains(err.Error(), "req_9") {
		t.Errorf("Error() = %q", err.Error())
	}
}

func TestAnErrorAnswerThatIsNotJSONStillReportsItsStatus(t *testing.T) {
	baseURL, _ := scriptedServer(t, reply{status: 502, headers: map[string]string{"Content-Type": "text/html"}, body: "<html>bad gateway</html>"})
	_, err := mustNew(t, "k", WithBaseURL(baseURL)).Esims.Search(context.Background(), nil)
	if e := asOpenRoamError(t, err); e.Status != 502 || e.Code != 0 {
		t.Errorf("error %+v", e)
	}
}

func TestNoAnswerIsAnOpenRoamErrorWithoutAStatusAndACancelStaysACancel(t *testing.T) {
	baseURL, _ := scriptedServer(t, reply{drop: true})
	_, err := mustNew(t, "k", WithBaseURL(baseURL)).Esims.Search(context.Background(), nil)
	if e := asOpenRoamError(t, err); e.Status != 0 || e.Unwrap() == nil {
		t.Errorf("error %+v, cause %v", e, e.Unwrap())
	}

	baseURL, _ = scriptedServer(t, okReply)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = mustNew(t, "k", WithBaseURL(baseURL)).Esims.Search(ctx, nil)
	if !errors.Is(err, context.Canceled) || asOpenRoamError(t, err).Status != 0 {
		t.Errorf("a cancelled call: %v", err)
	}

	slow := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
		case <-time.After(5 * time.Second):
		}
	}))
	defer slow.Close()
	ctx, cancel = context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	_, err = mustNew(t, "k", WithBaseURL(slow.URL), WithMaxRetries(2)).Esims.Search(ctx, nil)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("a call past its deadline: %v", err)
	}
}

func TestACtxThatEndsWhileAnErrorBodyIsReadStaysDetectable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(409)
		_, _ = w.Write([]byte(`{"code":`))
		w.(http.Flusher).Flush()
		<-r.Context().Done()
	}))
	defer server.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	_, err := mustNew(t, "k", WithBaseURL(server.URL)).Esims.Quote(ctx, esimOrder(t))
	if e := asOpenRoamError(t, err); e.Status != 409 || !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("error %+v, cause %v", e, e.Unwrap())
	}
}

func TestRetriesAreOffUntilWithMaxRetriesThenAnOrderIsSentAgainWithTheSameKeyAndBody(t *testing.T) {
	baseURL, requests := scriptedServer(t, reply{status: 503}, okReply)
	_, err := mustNew(t, "k", WithBaseURL(baseURL)).Esims.Order(context.Background(), esimOrder(t))
	if e := asOpenRoamError(t, err); e.Status != 503 || len(requests()) != 1 {
		t.Errorf("error %+v after %d requests", e, len(requests()))
	}

	baseURL, requests = scriptedServer(t, reply{status: 503, headers: map[string]string{"Retry-After": "0"}}, okReply)
	if _, err := mustNew(t, "k", WithBaseURL(baseURL), WithMaxRetries(2)).Esims.Order(context.Background(), esimOrder(t)); err != nil {
		t.Fatal(err)
	}
	got := requests()
	if len(got) != 2 || got[0].header.Get("Idempotency-Key") != got[1].header.Get("Idempotency-Key") || got[0].body != got[1].body {
		t.Errorf("requests %+v", got)
	}
}

func TestA202ThatSaysStillRunningIsAnErrorNotAnOrder(t *testing.T) {
	running := `{"code":12002,"name":"IDEMPOTENCY_KEY_IN_PROGRESS","message":"Still running.","requestId":"req_1"}`
	baseURL, _ := scriptedServer(t, reply{status: 202, body: running})
	order, err := mustNew(t, "k", WithBaseURL(baseURL)).Esims.Order(context.Background(), esimOrder(t))
	if e := asOpenRoamError(t, err); order != nil || e.Status != 202 || e.Code != 12002 || e.Name != "IDEMPOTENCY_KEY_IN_PROGRESS" {
		t.Errorf("error %+v, order %v", e, order)
	}

	baseURL, _ = scriptedServer(t, reply{status: 202, body: okReply.body})
	found, err := mustNew(t, "k", WithBaseURL(baseURL)).Esims.Search(context.Background(), nil)
	if err != nil || found == nil {
		t.Errorf("another 202: %v", err)
	}
}

func TestWithRetriesOnA12002AnswerIsRetriedAndTheFinalAnswerWins(t *testing.T) {
	baseURL, requests := scriptedServer(t, reply{status: 202, body: `{"code":12002}`, headers: map[string]string{"Retry-After": "0"}}, okReply)
	if _, err := mustNew(t, "k", WithBaseURL(baseURL), WithMaxRetries(2)).Esims.Order(context.Background(), esimOrder(t)); err != nil {
		t.Fatal(err)
	}
	if len(requests()) != 2 {
		t.Errorf("%d requests", len(requests()))
	}
}

func TestAMistakeInTheCallItselfIsAlsoAnOpenRoamErrorWithTheOriginalAsCause(t *testing.T) {
	baseURL, requests := scriptedServer(t, okReply)
	_, err := mustNew(t, "k", WithBaseURL(baseURL)).Airtime.LookupNumber(context.Background(), nil) // phoneNumber is required
	if e := asOpenRoamError(t, err); e.Status != 0 || e.Unwrap() == nil || len(requests()) != 0 {
		t.Errorf("error %+v, cause %v", e, e.Unwrap())
	}
}

func TestEveryAPIRouteHasItsMethodWithTheAgreedNames(t *testing.T) {
	methods := clientMethods(t)
	if len(methods) != 75 {
		t.Errorf("%d methods, want 75: %v", len(methods), methods)
	}
	for _, name := range []string{"Esims.Search", "Esims.Quote", "Esims.Order", "Esims.ListDailyUsage", "Esims.ListCatalog",
		"Esims.ListCompatibleDevices", "Esims.GetProfile", "OrderItems.Get", "OrderItems.Update", "OrderItems.Cancel", "OrderItems.ConfirmAmendment", "OrderItems.ManualAmendments.Cancel", "Airtime.LookupNumber", "Utilities.GetBiller", "GiftCards.GetGiftCard", "Hotels.GetHotel",
		"Hotels.ExportCatalog", "OrderItems.ManualAmendments.Approve", "Flights.GetSeatMap", "Orders.Order", "Orders.GetCancellationPreview",
		"OrderItems.GetCancellationPreview", "Wallets.SetDefault", "Countries.List", "MobileNetworks.List", "Currencies.List"} {
		if !slices.Contains(methods, name) {
			t.Errorf("%s is missing", name)
		}
	}
}

// The quick start of README.md compiles.
func TestTheReadmeQuickStartCompiles(t *testing.T) {
	_ = func(ctx context.Context) error {
		client, err := New("", WithEnvironment(Sandbox))
		if err != nil {
			return err
		}
		if _, err := client.Esims.Search(ctx, &EsimsSearchQuery{CountryCode: PtrString("JP")}); err != nil {
			var apiErr *OpenRoamError
			if errors.As(err, &apiErr) && apiErr.Code == 11002 {
				return nil
			}
			return err
		}
		items := []EsimOrderRequestItem{*NewEsimOrderRequestItem("esim_1")}
		body := NewEsimOrderRequest(items)
		body.SetCurrency("USD")
		if _, err := client.Esims.Quote(ctx, body); err != nil {
			return err
		}
		if _, err := client.Esims.Order(ctx, body, WithIdempotencyKey(NewIdempotencyKey())); err != nil {
			return err
		}
		if _, err := client.Esims.GetUsage(ctx, "iccid"); err != nil {
			return err
		}
		catalog, err := client.Esims.ListCatalog(ctx, &EsimsListCatalogQuery{Limit: PtrInt32(50)})
		if err != nil {
			return err
		}
		_, err = client.Esims.ListCatalog(ctx, &EsimsListCatalogQuery{Cursor: catalog.NextCursor.Get()})
		return err
	}
	_, _ = New("key",
		WithEnvironment(Sandbox),
		WithMaxRetries(3),
		WithHTTPClient(&http.Client{Timeout: 30 * time.Second}),
	)
}
