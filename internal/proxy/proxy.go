package proxy

import (
	"errors"
	"log/slog"
	"net/http"
	"net/textproto"

	"github.com/efixler/headless"
)

var (
	copyHeaders = []string{
		textproto.CanonicalMIMEHeaderKey("User-Agent"),
	}
)

func New(b headless.Browser) (http.HandlerFunc, error) {

	p := func(w http.ResponseWriter, req *http.Request) {
		slog.Debug("headless proxy request", "remote", req.RemoteAddr, "method", req.Method, "url", req.URL, "host", req.Host, "header", req.Header)
		headers := http.Header(make(map[string][]string))

		for _, headerName := range copyHeaders {
			if hval := req.Header.Get(headerName); hval != "" {
				headers.Set(headerName, hval)
			}
		}

		content, err := b.HTMLContent(req.URL.String(), headers)
		if err != nil {
			var httpErr *headless.HTTPError
			if errors.As(err, &httpErr) {
				http.Error(w, httpErr.Error(), httpErr.StatusCode)
			} else {
				http.Error(w, err.Error(), http.StatusBadGateway)
			}
			return
		}
		w.Write([]byte(content))
	}
	return p, nil
}
