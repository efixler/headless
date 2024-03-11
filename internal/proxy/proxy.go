package proxy

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/textproto"
	"strings"

	"github.com/efixler/headless"
	"github.com/efixler/headless/request"
)

type handlerMode int

const (
	AsProxy handlerMode = iota
	AsPostHandler
)

var (
	copyHeaders = []string{
		textproto.CanonicalMIMEHeaderKey("User-Agent"),
	}
)

type requestParser func(req *http.Request) (*request.Payload, error)

func New(b headless.TabFactory, mode handlerMode) (http.HandlerFunc, error) {
	var rp requestParser
	switch mode {
	case AsProxy:
		rp = parseProxyPayload
	case AsPostHandler:
		rp = parsePostPayload
	}

	p := func(w http.ResponseWriter, req *http.Request) {
		slog.Debug("headless proxy request", "remote", req.RemoteAddr, "method", req.Method, "url", req.URL, "host", req.Host, "header", req.Header)
		payload, err := rp(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		target, err := b.AcquireTab()
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		passHeaders := make(http.Header)
		for k, v := range payload.Headers {
			passHeaders.Set(k, v)
		}

		content, err := target.HTMLContent(payload.URL, passHeaders)
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

func parseProxyPayload(req *http.Request) (*request.Payload, error) {
	payload := &request.Payload{}
	payload.Headers = make(map[string]string)
	for _, headerName := range copyHeaders {
		if hval := req.Header.Get(headerName); hval != "" {
			payload.Headers[headerName] = hval
		}
	}
	payload.URL = req.URL.String()
	return payload, nil
}

func parsePostPayload(req *http.Request) (*request.Payload, error) {
	cType := strings.Split(req.Header.Get("Content-Type"), ";")[0]
	if cType != "application/json" {
		return nil, errors.New("Content-Type must be application/json")
	}
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	var payload request.Payload
	err := decoder.Decode(&payload)
	if err != nil {
		return nil, err
	}
	return &payload, nil
}
