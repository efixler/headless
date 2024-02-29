package browser

import (
	"context"
	"log/slog"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/efixler/headless/ua"
)

type Browser struct {
	ctx     context.Context
	Cancel  context.CancelFunc
	timeout time.Duration
}

func NewBrowser(ctx context.Context) *Browser {
	allocCtx, cancel := chromedp.NewExecAllocator(ctx, getAllocatorOptions()...)

	return &Browser{ctx: allocCtx, Cancel: cancel, timeout: 30 * time.Second}
}

func (b *Browser) HTMLContent(url string) (string, error) {
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

func getAllocatorOptions() []chromedp.ExecAllocatorOption {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		// chromedp.Flag("headless", false),
		// chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		// chromedp.Flag("mute-audio", true), // included in Headless
		// chromedp.Flag("remote-debugging-address", "127.0.0.1"),
		// chromedp.Flag("remote-debugging-port", fmt.Sprintf("%d", 9222)),
		// chromedp.WindowSize(1920, 1080),
		// chromedp.DisableGPU, // supposedly unnecessary
		chromedp.Headless,
		// chromedp.NoSandbox, TODO: Figure out what is better here in headless mode
		chromedp.UserAgent(ua.Firefox88),
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		// chromedp.IgnoreCertErrors, // check this when using proxies
		chromedp.WindowSize(1366, 768),
	)

	return opts
}
