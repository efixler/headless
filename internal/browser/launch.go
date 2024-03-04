package browser

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
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

type Chrome struct {
	ctx        context.Context
	Cancel     context.CancelFunc
	tabTimeout time.Duration
	sem        *semaphore.Weighted
	config     *config
}

type browserFunc func(url string, headers http.Header) (string, error)

func (f browserFunc) HTMLContent(url string, headers http.Header) (string, error) {
	return f(url, headers)
}

func NewChrome(ctx context.Context, options ...option) *Chrome {
	b := &Chrome{tabTimeout: 5 * time.Second}
	b.applyOptions(options)
	b.ctx, b.Cancel = chromedp.NewExecAllocator(ctx, b.config.allocatorOptions...)
	return b
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
	f := func(url string, headers http.Header) (string, error) {
		defer b.sem.Release(1)
		return b.HTMLContent(url, headers)
	}
	return browserFunc(f), nil
}

func (b *Chrome) HTMLContent(url string, headers http.Header) (string, error) {
	ctx, cancel := chromedp.NewContext(b.ctx)
	defer cancel()

	var html string
	var statusCode int

	listenCtx, cancelListen := context.WithCancel(ctx)
	chromedp.ListenTarget(listenCtx, func(ev interface{}) {
		if res, ok := ev.(*network.EventResponseReceived); ok {
			// see https://chromedevtools.github.io/devtools-protocol/tot/Network/#type-Response
			// eventChan <- res
			statusCode = int(res.Response.Status)
			slog.Info("Received response",
				"url", res.Response.URL,
				"type", res.Type,
				"status", statusCode,
				"statusText", res.Response.StatusText,
			)
			cancelListen()
		}
	})
	slog.Debug("Listening for response", "url", url)

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(1*time.Second),
		chromedp.WaitReady("body"),
		chromedp.OuterHTML("html", &html),
	)

	slog.Debug("Run completed HTML content", "url", url)
	if err != nil {
		// see https://github.com/chromedp/chromedp/blob/ebf842c7bc28db77d0bf4d757f5948d769d0866f/nav.go#L26
		// bad domain = page load error net::ERR_NAME_NOT_RESOLVED
		slog.Error("Error getting HTML content", "url", url, "err", err)
	}

	return html, err
}

// func getAllocatorOptions() []chromedp.ExecAllocatorOption {
// 	opts := append(chromedp.DefaultExecAllocatorOptions[:],
// 		// chromedp.Flag("headless", false),
// 		// chromedp.Flag("disable-gpu", true),
// 		chromedp.Flag("disable-dev-shm-usage", true),
// 		chromedp.Flag("blink-settings", "imagesEnabled=false"),
// 		// chromedp.Flag("mute-audio", true), // included in Headless
// 		// chromedp.Flag("remote-debugging-address", "127.0.0.1"),
// 		// chromedp.Flag("remote-debugging-port", fmt.Sprintf("%d", 9222)),
// 		// chromedp.WindowSize(1920, 1080),
// 		// chromedp.DisableGPU, // supposedly unnecessary
// 		chromedp.Headless,
// 		// chromedp.NoSandbox, TODO: Figure out what is better here in headless mode
// 		chromedp.UserAgent(ua.Firefox88),
// 		chromedp.NoFirstRun,
// 		chromedp.NoDefaultBrowserCheck,
// 		// chromedp.IgnoreCertErrors, // check this when using proxies
// 		chromedp.WindowSize(1366, 768),
// 	)

// 	return opts
// }
