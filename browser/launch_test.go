package browser

import (
	"context"
	"testing"
)

func TestGetErrorsOnInvalidURL(t *testing.T) {
	c := NewChrome(context.Background(), MaxTabs(1))
	b, err := c.AcquireTab()
	if err != nil {
		t.Fatalf("AcquireTab failed: %v", err)
	}
	_, err = b.Get("\\xyz::invalid.url", nil)
	if err == nil {
		t.Error("Get did not return an error")
	}
}

func TestExtractProtocolVersion(t *testing.T) {
	tests := []struct {
		p           string
		expectMajor int
		expectMinor int
	}{
		{
			p:           "HTTP/1.0",
			expectMajor: 1,
			expectMinor: 0,
		},
		{
			p:           "HTTP/1.1",
			expectMajor: 1,
			expectMinor: 1,
		},
		{
			p:           "http/1.1",
			expectMajor: 1,
			expectMinor: 1,
		},
		{
			p:           "HTTP/1",
			expectMajor: 1,
			expectMinor: 0,
		},
		{
			p:           "H2",
			expectMajor: 2,
			expectMinor: 0,
		},
		{
			p:           "",
			expectMajor: 1,
			expectMinor: 0,
		},
	}
	for _, tt := range tests {
		major, minor := extractHTTPVersion(tt.p)
		if major != tt.expectMajor {
			t.Errorf("[%s] expected major version %d, got %d", tt.p, tt.expectMajor, major)
		}
		if minor != tt.expectMinor {
			t.Errorf("[%s] expected minor version %d, got %d", tt.p, tt.expectMinor, minor)
		}
	}
}
