package proxy

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	nurl "net/url"
	"testing"
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
	proxyHandler, err := New(tp.mockBrowser)
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
