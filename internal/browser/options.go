package browser

import (
	"time"

	"github.com/chromedp/chromedp"
	"github.com/efixler/headless/ua"
	"golang.org/x/sync/semaphore"
)

type config struct {
	allocatorOptions []chromedp.ExecAllocatorOption
	userAgent        string
	windowSize       [2]int
}

type option func(*Chrome) error

func CDPOptions(cdps ...chromedp.ExecAllocatorOption) option {
	return func(b *Chrome) error {
		b.config.allocatorOptions = append(b.config.allocatorOptions, cdps...)
		return nil
	}
}

func Headless(h bool) option {
	return func(b *Chrome) error {
		if h {
			b.config.allocatorOptions = append(b.config.allocatorOptions, chromedp.Headless)
		}
		return nil
	}
}

func UserAgentIfNotEmpty(ua string) option {
	return func(b *Chrome) error {
		if ua != "" {
			b.config.userAgent = ua
		}
		return nil
	}
}

func MaxTabs(n int) option {
	return func(b *Chrome) error {
		b.sem = semaphore.NewWeighted(int64(n))
		return nil
	}
}

func TabAcquireTimeout(d time.Duration) option {
	return func(b *Chrome) error {
		b.tabTimeout = d
		return nil
	}
}

func WindowSize(w, h int) option {
	return func(b *Chrome) error {
		b.config.windowSize = [2]int{w, h}
		return nil
	}
}

func getDefaults() config {
	return config{
		allocatorOptions: []chromedp.ExecAllocatorOption{
			chromedp.Flag("disable-dev-shm-usage", true),
			chromedp.Flag("blink-settings", "imagesEnabled=false"),
			chromedp.NoFirstRun,
			chromedp.NoDefaultBrowserCheck,
			// chromedp.NoSandbox, TODO: Figure out what is better here in headless mode
			// chromedp.Flag("disable-software-rasterizer", true),
			// chromedp.IgnoreCertErrors, // check this when using proxies
		},
	}
}

func (b *Chrome) applyOptions(opts []option) error {
	c := getDefaults()
	b.config = &c
	for _, opt := range opts {
		if err := opt(b); err != nil {
			return err
		}
	}
	if b.config.userAgent == "" {
		b.config.userAgent = ua.Firefox88
	}

	c.allocatorOptions = append(c.allocatorOptions,
		chromedp.UserAgent(b.config.userAgent),
		chromedp.WindowSize(b.config.windowSize[0], b.config.windowSize[1]),
	)

	return nil
}
