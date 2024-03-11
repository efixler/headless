package proxy

import (
	"net/http"

	"github.com/efixler/headless"
)

func Service(c headless.TabFactory) (http.Handler, error) {
	headlessHandler, err := New(c, AsPostHandler)
	if err != nil {
		return nil, err
	}
	mux := http.NewServeMux()
	mux.Handle("POST /{$}", headlessHandler)
	return mux, nil
}

func HTTPProxy(c headless.TabFactory) (http.Handler, error) {
	return New(c, AsProxy)
}
