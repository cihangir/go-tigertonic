package tigertonic

import (
	"net/http"
	"net/url"
	"testing"
)

type TestResponse struct {
	ImportantInfo string `json:"important_info"`
}

// GET /baz
func get(u *url.URL, h http.Header, _ interface{}) (int, http.Header, *TestResponse, error) {
	return http.StatusOK, nil, &TestResponse{"i love you"}, nil
}

func TestCORSOPTIONS(t *testing.T) {
	mux := NewTrieServeMux()
	mux.Handle("GET", "/foo", NewCORSBuilder().AddAllowedOrigins("*").Build(Marshaled(get)))
	mux.Handle("GET", "/baz", NewCORSBuilder().AddAllowedOrigins("http://gooddomain.com").Build(Marshaled(get)))
	mux.Handle("GET", "/quux", NewCORSBuilder().AddAllowedHeaders("X-Pizza-Fax").Build(Marshaled(get)))

	w := &testResponseWriter{}
	r, _ := http.NewRequest("OPTIONS", "http://example.com/baz", nil)
	r.Header.Set(CORSRequestMethod, "GET")
	mux.ServeHTTP(w, r)
	if http.StatusOK != w.StatusCode {
		t.Fatal(w.StatusCode)
	}
	if "GET, HEAD, OPTIONS" != w.Header().Get(CORSAllowMethods) {
		t.Fatal(w.Header().Get("Allow"))
	}

	// requesting secured resource with invalid domain
	w = &testResponseWriter{}
	r, _ = http.NewRequest("OPTIONS", "http://example.com/baz", nil)
	r.Header.Set(CORSRequestOrigin, "http://baddomain.com")
	r.Header.Set(CORSRequestMethod, "GET")
	mux.ServeHTTP(w, r)
	if http.StatusOK != w.StatusCode {
		t.Fatal(w.StatusCode)
	}
	if "null" != w.Header().Get(CORSAllowOrigin) {
		t.Fatal(w.Header().Get(CORSAllowOrigin))
	}

	// requesting unsecured/wildcard resource with invalid domain
	w = &testResponseWriter{}
	r, _ = http.NewRequest("OPTIONS", "http://example.com/foo", nil)
	r.Header.Set(CORSRequestOrigin, "http://baddomain.com")
	r.Header.Set(CORSRequestMethod, "GET")
	mux.ServeHTTP(w, r)
	if http.StatusOK != w.StatusCode {
		t.Fatal(w.StatusCode)
	}
	if "*" != w.Header().Get(CORSAllowOrigin) {
		t.Fatal(w.Header().Get(CORSAllowOrigin))
	}

	// requesting secured resource with valid domain
	w = &testResponseWriter{}
	r, _ = http.NewRequest("OPTIONS", "http://example.com/baz", nil)
	r.Header.Set(CORSRequestOrigin, "http://gooddomain.com")
	r.Header.Set(CORSRequestMethod, "GET")
	mux.ServeHTTP(w, r)
	if http.StatusOK != w.StatusCode {
		t.Fatal(w.StatusCode)
	}
	if "http://gooddomain.com" != w.Header().Get(CORSAllowOrigin) {
		t.Fatal(w.Header().Get(CORSAllowOrigin))
	}

	// just requesting some headers, mane
	w = &testResponseWriter{}
	r, _ = http.NewRequest("OPTIONS", "http://example.com/quux", nil)
	r.Header.Set(CORSRequestMethod, "GET")
	r.Header.Add(CORSRequestHeaders, "X-Pizza-Fax")
	mux.ServeHTTP(w, r)
	if http.StatusOK != w.StatusCode {
		t.Fatal(w.StatusCode)
	}
	t.Log(w.Header())
	if "X-Pizza-Fax" != w.Header().Get(CORSAllowHeaders) {
		t.Fatalf("Headers received missing pizza fax! %s", w.Header())
	}
}

func TestCORSOrigin(t *testing.T) {
	mux := NewTrieServeMux()
	mux.Handle("GET", "/foo", NewCORSBuilder().AddAllowedOrigins("*").Build(Marshaled(get)))
	mux.Handle("GET", "/baz", NewCORSBuilder().AddAllowedOrigins("http://gooddomain.com").Build(Marshaled(get)))

	// wildcard
	w := &testResponseWriter{}
	r, _ := http.NewRequest("GET", "http://example.com/foo", nil)
	r.Header.Set("Accept", "application/json")
	r.Header.Set(CORSRequestOrigin, "http://gooddomain.com")
	mux.ServeHTTP(w, r)
	if http.StatusOK != w.StatusCode {
		t.Fatal(w.StatusCode)
	}
	if "*" != w.Header().Get(CORSAllowOrigin) {
		t.Fatal(w.Header().Get(CORSAllowOrigin))
	}

	// specific
	w = &testResponseWriter{}
	r, _ = http.NewRequest("GET", "http://example.com/baz", nil)
	r.Header.Set("Accept", "application/json")
	r.Header.Set(CORSRequestOrigin, "http://gooddomain.com")
	mux.ServeHTTP(w, r)
	if http.StatusOK != w.StatusCode {
		t.Fatal(w.StatusCode)
	}
	if "http://gooddomain.com" != w.Header().Get(CORSAllowOrigin) {
		t.Fatal(w.Header().Get(CORSAllowOrigin))
	}

}
func TestCORSHeader(t *testing.T) {
	mux := NewTrieServeMux()
	mux.Handle("GET", "/foo", NewCORSBuilder().Build(Marshaled(get)))
	mux.Handle("GET", "/baz", NewCORSBuilder().AddAllowedHeaders("X-Fancy-Header").Build(Marshaled(get)))

	// wildcard
	w := &testResponseWriter{}
	r, _ := http.NewRequest("GET", "http://example.com/foo", nil)
	r.Header.Set("Accept", "application/json")
	r.Header.Set(CORSRequestHeaders, "X-Header-You-Dont-Want")
	mux.ServeHTTP(w, r)
	if http.StatusOK != w.StatusCode {
		t.Fatal(w.StatusCode)
	}
	if "" != w.Header().Get(CORSAllowHeaders) {
		t.Fatal(w.Header().Get(CORSAllowOrigin))
	}

	// specific
	w = &testResponseWriter{}
	r, _ = http.NewRequest("GET", "http://example.com/baz", nil)
	r.Header.Set("Accept", "application/json")
	r.Header.Set(CORSRequestHeaders, "X-Fancy-Header")
	mux.ServeHTTP(w, r)
	if http.StatusOK != w.StatusCode {
		t.Fatal(w.StatusCode)
	}
	if "X-Fancy-Header" != w.Header().Get(CORSAllowHeaders) {
		t.Fatal(w.Header().Get(CORSAllowOrigin))
	}

}
