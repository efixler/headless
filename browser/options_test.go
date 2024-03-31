package browser

import (
	"context"
	"testing"

	"github.com/efixler/headless/ua"
)

func TestUserAgentIfNotEmpty(t *testing.T) {
	c, err := NewChrome(context.Background())
	if err != nil {
		t.Fatalf("NewChrome failed: %v", err)
	}
	if c.config.userAgent != "" {
		t.Errorf("UserAgent option, expected empty got %s", c.config.userAgent)
	}
	type data struct {
		name     string
		args     []string
		expected string
	}
	tests := []data{
		{"empty", []string{""}, ""},
		{"single", []string{"User agent"}, "User agent"},
		{"multiple, first empty", []string{"", "second", "third"}, "third"},
	}

	for _, test := range tests {
		opt := UserAgentIfNotEmpty(test.args...)
		opt(c)
		if c.config.userAgent != test.expected {
			t.Errorf("[%s] expected %s got %s", test.name, test.expected, c.config.userAgent)
		}
	}
	opt := AsFirefox()
	opt(c)
	if c.config.userAgent != ua.Firefox88 {
		t.Errorf("UserAgent option, expected empty got %s", c.config.userAgent)
	}

}
