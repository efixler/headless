package proxy

import (
	"net/http/httptest"
	"testing"
)

func TestServiceOnlyAllowsPostToHeadless(t *testing.T) {
	tf := &mockBrowser{}
	handler, err := Service(tf)
	if err != nil {
		t.Fatalf("Service() error: %v", err)
	}
	for _, method := range []string{"GET", "PUT", "DELETE", "PATCH", "OPTIONS", "CONNECT", "TRACE"} {
		req := httptest.NewRequest(method, "/", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		if w.Code != 405 {
			t.Errorf("Expected 405, got %d", w.Code)
		}
	}
}
