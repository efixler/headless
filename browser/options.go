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

type ChromeOption func(*Chrome) error

func CDPOptions(cdps ...chromedp.ExecAllocatorOption) ChromeOption {
	return func(b *Chrome) error {
		b.config.allocatorOptions = append(b.config.allocatorOptions, cdps...)
		return nil
	}
}

func Headless(h bool) ChromeOption {
	return func(b *Chrome) error {
		if h {
			b.config.allocatorOptions = append(b.config.allocatorOptions, chromedp.Headless)
		}
		return nil
	}
}

func AsFirefox() ChromeOption {
	return UserAgentIfNotEmpty(ua.Firefox88)
}

func UserAgentIfNotEmpty(uas ...string) ChromeOption {
	return func(b *Chrome) error {
		for _, agent := range uas {
			if agent == "" {
				continue
			}
			b.config.userAgent = agent
		}
		return nil
	}
}

func MaxTabs(n int) ChromeOption {
	return func(b *Chrome) error {
		b.sem = semaphore.NewWeighted(int64(n))
		return nil
	}
}

func TabAcquireTimeout(d time.Duration) ChromeOption {
	return func(b *Chrome) error {
		b.tabTimeout = d
		return nil
	}
}

func WindowSize(w, h int) ChromeOption {
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
			// chromedp.Flag("mute-audio", true), // included in Headless
		},
		windowSize: [2]int{1366, 768},
	}
}

func (b *Chrome) applyOptions(opts []ChromeOption) error {
	c := getDefaults()
	b.config = &c
	for _, opt := range opts {
		if err := opt(b); err != nil {
			return err
		}
	}

	if b.config.userAgent != "" {
		c.allocatorOptions = append(
			c.allocatorOptions,
			chromedp.UserAgent(b.config.userAgent),
		)
	}

	c.allocatorOptions = append(c.allocatorOptions,
		chromedp.WindowSize(b.config.windowSize[0], b.config.windowSize[1]),
	)

	return nil
}
