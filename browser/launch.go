package browser

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/efixler/headless"
	"golang.org/x/sync/semaphore"
)

var (
	// ErrMaxTabs is returned when the maximum number of tabs is reached
	ErrMaxTabs       = fmt.Errorf("maximum number of tabs reached")
	ErrMaxTabsNotSet = fmt.Errorf("maximum number of tabs not set")
)

func NewChrome(ctx context.Context, options ...option) *Chrome {
	b := &Chrome{tabTimeout: 5 * time.Second}
	b.applyOptions(options)
	b.ctx, b.Cancel = chromedp.NewExecAllocator(ctx, b.config.allocatorOptions...)
	return b
}

type Chrome struct {
	ctx        context.Context
	Cancel     context.CancelFunc
	tabTimeout time.Duration
	sem        *semaphore.Weighted
	config     *config
}

type browserFunc func(url string, headers http.Header) (*http.Response, error)

func (f browserFunc) Get(url string, headers http.Header) (*http.Response, error) {
	return f(url, headers)
}

func (b *Chrome) AcquireTab() (headless.Browser, error) {
	if b.sem == nil {
		return nil, ErrMaxTabsNotSet
	}
	tabWaitContext, cancel := context.WithTimeout(b.ctx, b.tabTimeout)
	defer cancel()
	if err := b.sem.Acquire(tabWaitContext, 1); err != nil {
		return nil, errors.Join(err, ErrMaxTabs)
	}

	f := func(url string, headers http.Header) (*http.Response, error) {
		defer b.sem.Release(1)
		return b.Get(url, headers)
	}
	return browserFunc(f), nil
}

func (b *Chrome) Get(url string, headers http.Header) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	ctx, cancel := chromedp.NewContext(b.ctx)
	defer cancel()

	var html string
	response := &http.Response{
		Header: http.Header{},
	}

	listenCtx, cancelListen := context.WithCancel(ctx)
	chromedp.ListenTarget(listenCtx, func(ev interface{}) {
		if res, ok := ev.(*network.EventResponseReceived); ok {
			// see https://chromedevtools.github.io/devtools-protocol/tot/Network/#type-Response
			// The first event should be a 'Document' type response with the status of the page load
			response.StatusCode = int(res.Response.Status)
			response.Status = fmt.Sprintf("%d %s", response.StatusCode, res.Response.StatusText)
			response.Proto = strings.ToUpper(res.Response.Protocol)
			response.ProtoMajor, response.ProtoMinor = extractHTTPVersion(res.Response.Protocol)
			slog.Debug("Received network event from page",
				"url", res.Response.URL,
				"type", res.Type,
				"status", response.Status,
				"statusText", res.Response.StatusText,
				"protocol", res.Response.Protocol,
			)
			headers := res.Response.Headers
			for k, v := range headers {
				switch http.CanonicalHeaderKey(k) {
				case "Content-Length":
					continue
				case "Content-Encoding":
					continue
				}
				response.Header.Add(k, fmt.Sprintf("%v", v))
				slog.Debug("Header", "key", k, "value", v)
			}
			cancelListen()
		}
	})
	slog.Debug("Navigating to:", "url", url)
	// TODO: add passHeaders to request
	err = chromedp.Run(ctx,
		chromedp.Navigate(request.URL.String()),
		chromedp.Sleep(1*time.Second),
		chromedp.WaitReady("body"),
		chromedp.OuterHTML("html", &html),
	)

	if err != nil {
		// see https://github.com/chromedp/chromedp/blob/ebf842c7bc28db77d0bf4d757f5948d769d0866f/nav.go#L26
		// bad domain = page load error net::ERR_NAME_NOT_RESOLVED
		slog.Error("Error getting HTML content", "url", url, "err", err)
	}
	response.ContentLength = int64(len(html))
	response.Body = io.NopCloser(strings.NewReader(html))
	response.Request = request
	return response, err
}

func extractHTTPVersion(protocol string) (major, minor int) {
	major = 1
	protocol = strings.ToUpper(protocol)
	n, _ := fmt.Sscanf(protocol, "HTTP/%d.%d", &major, &minor)
	if (n > 0) || (len(protocol) < 2) {
		return
	}
	if protocol[0:2] == "H2" {
		major = 2
	}
	return
}
