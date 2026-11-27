package openroam

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"slices"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"
)

// One row of testdata/routes.json, generated from the API description: a client method and the request it must make.
type routeRow struct {
	Method    string   `json:"method"`
	Verb      string   `json:"verb"`
	Route     string   `json:"route"`
	PathNames []string `json:"pathNames"`
	Query     []struct {
		Name     string `json:"name"`
		Required bool   `json:"required"`
		Kind     string `json:"kind"`
	} `json:"query"`
	Idempotent bool            `json:"idempotent"`
	Body       bool            `json:"body"`
	BodySample json.RawMessage `json:"bodySample"`
}

// The samples routes_calls_test.go sends as query parameters, and what the server must receive for each kind.
var (
	sampleString = "sample"
	sampleTime   = time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	sampleArray  = []string{"a", "b"}
	sampleSent   = map[string]string{"string": "sample", "number": "2", "boolean": "true", "date": "2026-01-02T03:04:05Z"}
)

// readSample reads a sample body into the model of the route, as a caller would hold it.
func readSample(t *testing.T, body []byte, into any) {
	t.Helper()
	if err := json.Unmarshal(body, into); err != nil {
		t.Fatalf("the sample body does not read into %T: %v", into, err)
	}
}

func loadRoutes(t *testing.T) []routeRow {
	t.Helper()
	data, err := os.ReadFile("testdata/routes.json")
	if err != nil {
		t.Fatal(err)
	}
	var routes []routeRow
	if err := json.Unmarshal(data, &routes); err != nil {
		t.Fatal(err)
	}
	return routes
}

type recorded struct {
	method string
	url    *url.URL
	header http.Header
	body   string
}

// recordRequests runs use with a client whose server records each request and answers with an empty object.
func recordRequests(t *testing.T, use func(c *OpenRoam)) []recorded {
	t.Helper()
	var mu sync.Mutex
	var requests []recorded
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		mu.Lock()
		requests = append(requests, recorded{r.Method, r.URL, r.Header.Clone(), string(body)})
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{}"))
	}))
	defer server.Close()
	c, err := New("key", WithBaseURL(server.URL))
	if err != nil {
		t.Fatal(err)
	}
	use(c)
	mu.Lock()
	defer mu.Unlock()
	return requests
}

func TestTheRouteTableCoversAll73RoutesEachWithItsOwnMethod(t *testing.T) {
	routes := loadRoutes(t)
	if len(routes) != 73 {
		t.Fatalf("%d routes, want 73", len(routes))
	}
	seen := map[string]bool{}
	for _, route := range routes {
		if seen[route.Method] {
			t.Errorf("%s twice", route.Method)
		}
		seen[route.Method] = true
		if routeCalls[route.Method] == nil {
			t.Errorf("no call for %s", route.Method)
		}
	}
	if len(routeCalls) != len(routes) {
		t.Errorf("%d calls for %d routes", len(routeCalls), len(routes))
	}
}

func TestEveryRouteMakesItsRequest(t *testing.T) {
	for _, route := range loadRoutes(t) {
		t.Run(route.Method, func(t *testing.T) {
			call := routeCalls[route.Method]
			if call == nil {
				t.Fatalf("no call for %s", route.Method)
			}
			ids := make([]string, len(route.PathNames))
			for i, name := range route.PathNames {
				ids[i] = name + "-value"
			}
			requests := recordRequests(t, func(c *OpenRoam) {
				// The answer is an empty object, which some response types refuse to read: this test is about the request.
				_ = call(t, context.Background(), c, ids, route.BodySample)
			})
			if len(requests) != 1 {
				t.Fatalf("%d requests, want exactly one", len(requests))
			}
			request := requests[0]
			if request.method != route.Verb {
				t.Errorf("verb %s, want %s", request.method, route.Verb)
			}
			want := route.Route
			for i, name := range route.PathNames {
				want = strings.Replace(want, "{"+name+"}", url.PathEscape(ids[i]), 1)
			}
			if request.url.EscapedPath() != want {
				t.Errorf("path %s, want %s", request.url.EscapedPath(), want)
			}
			if got := request.header.Get("X-Partner-Key"); got != "key" {
				t.Errorf("X-Partner-Key %q", got)
			}

			query := request.url.Query()
			for _, q := range route.Query {
				if !query.Has(q.Name) {
					t.Errorf("query parameter %s is not sent", q.Name)
					continue
				}
				if q.Kind == "array" {
					if got := query[q.Name]; !slices.Equal(got, sampleArray) {
						t.Errorf("query parameter %s = %v, want %v", q.Name, got, sampleArray)
					}
				} else if got := query.Get(q.Name); got != sampleSent[q.Kind] {
					t.Errorf("query parameter %s = %q, want %q", q.Name, got, sampleSent[q.Kind])
				}
			}
			if len(route.Query) == 0 && request.url.RawQuery != "" {
				t.Errorf("unexpected query %q", request.url.RawQuery)
			}

			if route.Body {
				// The body itself is sent, flat: its own top-level fields, not wrapped in an object named after its type.
				var sent, sample map[string]any
				// JSON null reads into a nil map without an error, so it is refused by name.
				if err := json.Unmarshal([]byte(request.body), &sent); err != nil || sent == nil {
					t.Fatalf("the body is not a JSON object: %q", request.body)
				}
				if err := json.Unmarshal(route.BodySample, &sample); err != nil {
					t.Fatal(err)
				}
				if got, want := keys(sent), keys(sample); !slices.Equal(got, want) {
					t.Errorf("body fields %v, want %v", got, want)
				}
				if !strings.Contains(request.header.Get("Content-Type"), "application/json") {
					t.Errorf("Content-Type %q", request.header.Get("Content-Type"))
				}
			} else if request.body != "" {
				t.Errorf("unexpected body %q", request.body)
			}

			key := request.header.Get("Idempotency-Key")
			if route.Idempotent != (key != "") {
				t.Errorf("Idempotency-Key %q on a route where it is needed: %v", key, route.Idempotent)
			}
			if route.Idempotent && !uuidV4.MatchString(key) {
				t.Errorf("Idempotency-Key %q is not a UUID v4", key)
			}

			if route.Idempotent {
				own := recordRequests(t, func(c *OpenRoam) {
					_ = call(t, context.Background(), c, ids, route.BodySample, WithIdempotencyKey("my-key-42"))
				})
				if len(own) != 1 || own[0].header.Get("Idempotency-Key") != "my-key-42" {
					t.Errorf("your own key is not sent: %v", own)
				}
			}
		})
	}
}

// Every method of the route table exists on the client under its agreed name, in PascalCase.
func TestEveryRouteIsAMethodOfTheClient(t *testing.T) {
	methods := clientMethods(t)
	for _, route := range loadRoutes(t) {
		parts := strings.Split(route.Method, ".")
		for i, part := range parts {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
		if !slices.Contains(methods, strings.Join(parts, ".")) {
			t.Errorf("%s is not a method of the client", strings.Join(parts, "."))
		}
	}
}

// clientMethods lists every method of every resource of the client, for example "OrderItems.GetChange".
func clientMethods(t *testing.T) []string {
	t.Helper()
	c, err := New("key")
	if err != nil {
		t.Fatal(err)
	}
	var methods []string
	var walk func(value reflect.Value, prefix string)
	walk = func(value reflect.Value, prefix string) {
		for i := 0; i < value.NumMethod(); i++ {
			methods = append(methods, prefix+value.Type().Method(i).Name)
		}
		elem := value.Elem()
		for i := 0; i < elem.NumField(); i++ {
			field := elem.Type().Field(i)
			if field.IsExported() && field.Type.Kind() == reflect.Pointer && strings.HasSuffix(field.Type.Elem().Name(), "Resource") {
				walk(elem.Field(i), prefix+field.Name+".")
			}
		}
	}
	walk(reflect.ValueOf(c), "")
	sort.Strings(methods)
	return methods
}

func keys(m map[string]any) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
