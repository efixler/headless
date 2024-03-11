package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	nurl "net/url"
	"strings"
	"testing"

	"github.com/efixler/headless"
	"github.com/efixler/headless/request"
)

type testProxy struct {
	targetServer *httptest.Server
	proxyServer  *httptest.Server
	client       *http.Client
	mockBrowser  *mockBrowser
}

type mockBrowser struct {
	headers http.Header
	url     string
}

func (b *mockBrowser) AcquireTab() (headless.Browser, error) {
	return b, nil
}

func (b *mockBrowser) HTMLContent(url string, headers http.Header) (string, error) {
	b.url = url
	b.headers = headers
	html := fmt.Sprintf(`<html><head><title>%s</title></head><body>%s</body></html>`, url, url)
	return html, nil
}

func newTestProxy(t *testing.T) *testProxy {
	tp := &testProxy{}
	targetHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// the request should never actually get here because it is satisfied by the proxy
		t.Fatalf("Unexpected pass-through url: %q", r.URL)
	})
	tp.targetServer = httptest.NewServer(targetHandler)
	tp.mockBrowser = &mockBrowser{}
	proxyHandler, err := New(tp.mockBrowser, AsProxy)
	if err != nil {
		t.Fatalf("can't initialize proxy handler %v", err)
	}
	tp.proxyServer = httptest.NewServer(proxyHandler)
	proxyURL, _ := nurl.Parse(tp.proxyServer.URL)
	tp.client = tp.targetServer.Client()
	tp.client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	return tp
}

func (tp *testProxy) close() {
	tp.proxyServer.Close()
	tp.targetServer.Close()

}

// This test is HTTP-specific and mainly tests that the url gets through
// correctly. It also demonstrates how to set up a client to use the proxy.
func TestBasicProxyRequest(t *testing.T) {
	tp := newTestProxy(t)
	defer tp.close()
	type data struct {
		name         string
		url          string
		expectStatus int
	}
	tests := []data{
		{"basic", "http://www.foobar.com/today/index.html", 200},
	}
	for _, test := range tests {
		r, err := tp.client.Get(test.url)
		if err != nil {
			t.Fatalf("[%s] - unexpected error %v", test.name, err)
		}
		defer r.Body.Close()
		if r.StatusCode != test.expectStatus {
			t.Errorf("[%s] - expected status %d, got %d", test.name, test.expectStatus, r.StatusCode)
		}
		if tp.mockBrowser.url != test.url {
			t.Errorf("[%s] - expected url %s, got %s", test.name, test.url, tp.mockBrowser.url)
		}
	}
}

func TestProxyAsPostHandler(t *testing.T) {
	tests := []struct {
		name          string
		url           string
		header        map[string]string
		expectStatus  int
		expectHeaders map[string]string
	}{
		{
			"basic",
			"http://www.foobar.com/today/index.html",
			map[string]string{"User-Agent": "fooagent", "Accept": "text/html"},
			200,
			map[string]string{"User-Agent": "fooagent"},
		},
		{
			"nil headers",
			"http://www.foobar.com/today/index.html",
			nil,
			200,
			map[string]string{},
		},
	}
	mockBrowser := mockBrowser{}
	headlessHandler, err := New(&mockBrowser, AsPostHandler)
	if err != nil {
		t.Fatalf("can't initialize proxy handler %v", err)
	}
	for _, test := range tests {
		payload := &request.Payload{URL: test.url, Headers: test.header}
		payloadJSON, err := json.Marshal(payload)
		body := bytes.NewReader(payloadJSON)
		if err != nil {
			t.Fatalf("[%s] - unexpected error marshaling json %v", test.name, err)
		}

		req := httptest.NewRequest("POST", "/", body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		headlessHandler(w, req)
		resp := w.Result()
		if resp.StatusCode != test.expectStatus {
			t.Errorf("[%s] - expected status %d, got %d", test.name, test.expectStatus, w.Code)
		}
		for k, v := range test.expectHeaders {
			if mockBrowser.headers.Get(k) != v {
				t.Errorf("[%s] - expected header %s: %v, got %v", test.name, k, v, mockBrowser.headers.Get(k))
			}
		}
	}
}

func TestPayloadBadContentType(t *testing.T) {
	mockBrowser := mockBrowser{}
	headlessHandler, err := New(&mockBrowser, AsPostHandler)
	if err != nil {
		t.Fatalf("can't initialize proxy handler %v", err)
	}
	payload := &request.Payload{URL: "foo", Headers: map[string]string{"Content-Type": "application/json"}}
	payloadJSON, _ := json.Marshal(payload)
	body := bytes.NewReader(payloadJSON)
	req := httptest.NewRequest("POST", "/", body)
	req.Header.Set("Content-Type", "text/plain")
	w := httptest.NewRecorder()
	headlessHandler(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestPayloadDisallowedFields(t *testing.T) {
	mockBrowser := mockBrowser{}
	headlessHandler, err := New(&mockBrowser, AsPostHandler)
	if err != nil {
		t.Fatalf("can't initialize proxy handler %v", err)
	}
	payloadJSON := `{"url": "foo", "headers": {"Content-Type": "application/json"}, "foo": "bar"}`
	body := strings.NewReader(payloadJSON)
	req := httptest.NewRequest("POST", "/", body)
	w := httptest.NewRecorder()
	headlessHandler(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}
