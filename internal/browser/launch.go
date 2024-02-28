package browser

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/efixler/headless/ua"
)

type Browser struct {
	ctx    context.Context
	Cancel context.CancelFunc
}

func NewBrowser(ctx context.Context) *Browser {
	allocCtx, cancel := chromedp.NewExecAllocator(ctx, getAllocatorOptions()...)

	return &Browser{ctx: allocCtx, Cancel: cancel}
}

func (b *Browser) HTMLContent(url string) (string, error) {
	ctx, cancel := chromedp.NewContext(b.ctx)
	defer cancel()

	var htmlContent string
	if err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
		chromedp.Sleep(2*time.Second),
		chromedp.OuterHTML("html", &htmlContent),
	); err != nil {
		return "", err
	}
	return htmlContent, nil
}

func getAllocatorOptions() []chromedp.ExecAllocatorOption {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		// chromedp.Flag("headless", true),
		// chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-dev-shm-usage", true),
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
