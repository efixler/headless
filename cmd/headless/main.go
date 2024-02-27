package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/efixler/headless/internal/browser"
	"github.com/efixler/headless/internal/cmd"
)

var (
	flags        = flag.NewFlagSet("what is the name for?", flag.ExitOnError)
	browserFlags = cmd.HeadlessBrowserFlags("HEADLESS_", flags)
)

func main() {
	fmt.Printf("Address: %s\n", browserFlags.Address.Get())
	fmt.Printf("Port: %d\n", browserFlags.Port.Get())
	fmt.Printf("URL: %v\n", flags.Args())
	if len(flags.Args()) == 0 {
		flags.Usage()
		os.Exit(1)
	}
	url := flags.Args()[0]

	b := browser.NewBrowser(context.Background())
	defer b.Cancel()

	// content, err := GetHtmlContent(url)
	content, err := b.HTMLContent(url)
	if err != nil {
		slog.Error("Error getting HTML content", "url", url, "err", err)
		os.Exit(1)
	}
	fmt.Println(content)
}

// func GetHtmlContent(url string) (string, error) {
// 	// ctx, cancel := chromedp.NewContext(context.Background())
// 	// defer cancel()

// 	// ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
// 	// defer cancel()

// 	// opts := append(chromedp.DefaultExecAllocatorOptions[:],
// 	// 	chromedp.Flag("headless", false),
// 	// 	chromedp.Flag("disable-gpu", true),
// 	// 	chromedp.Flag("remote-debugging-address", browserFlags.Address.Get()),
// 	// 	chromedp.Flag("remote-debugging-port", fmt.Sprintf("%d", browserFlags.Port.Get())),
// 	// )
// 	// allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)

// 	// allocCtx, cancel := chromedp.NewRemoteAllocator(context.Background(),
// 	// 	fmt.Sprintf("http://%s:%d", browserFlags.Address.Get(), browserFlags.Port.Get()),
// 	// )

// 	b := browser.NewBrowser(context.Background())

// 	defer b.Cancel()
// 	ctx, cancel := chromedp.NewContext(allocCtx)
// 	defer cancel()

// 	if err := chromedp.Run(ctx,
// 		chromedp.Navigate(url),
// 		chromedp.WaitVisible("body"),
// 	); err != nil {
// 		return "", err
// 	}

// 	var htmlContent string
// 	if err := chromedp.Run(ctx, chromedp.OuterHTML("html", &htmlContent)); err != nil {
// 		return "", err
// 	}

// 	return htmlContent, nil
// }

func init() {
	flags.Usage = usage
	flags.Parse(os.Args[1:])
}

func usage() {
	fmt.Println(`Usage: 
	headless [flags] :url
 
  -h	
  	Show this help message`)

	flags.PrintDefaults()
}
