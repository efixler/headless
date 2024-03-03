package headless

import (
	"fmt"
	"net/http"
	"strings"
)

type Browser interface {
	HTMLContent(url string, headers http.Header) (string, error)
}

type HTTPError struct {
	StatusCode int
	Status     string
	Message    string
}

func (e HTTPError) Error() string {
	if e.Status == "" {
		e.Status = http.StatusText(e.StatusCode)
	}
	return strings.TrimSpace(fmt.Sprintf("%s %s", e.Status, e.Message))
}

func (e HTTPError) String() string {
	if e.Status == "" {
		e.Status = http.StatusText(e.StatusCode)
	}
	return strings.TrimSpace(fmt.Sprintf("HTTP error %d:%s %s", e.StatusCode, e.Status, e.Message))
}
