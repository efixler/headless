package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/efixler/envflags"
	"github.com/efixler/headless/internal/browser"
	"github.com/efixler/headless/internal/cmd"
)

var (
	flags        = flag.NewFlagSet("headless", flag.ExitOnError)
	browserFlags *cmd.HeadlessBrowserSpec
	headless     bool
)

func main() {
	fmt.Fprintf(os.Stderr, "Address: %s\n", browserFlags.Address.Get())
	fmt.Fprintf(os.Stderr, "Port: %d\n", browserFlags.Port.Get())
	fmt.Fprintf(os.Stderr, "URL: %v\n", flags.Args())
	if len(flags.Args()) == 0 {
		flags.Usage()
		os.Exit(1)
	}
	url := flags.Args()[0]

	b := browser.NewChrome(
		context.Background(),
		browser.Headless(headless),
		browser.MaxTabs(1),
	)
	defer b.Cancel()
	tab, err := b.AcquireTab()
	if err != nil {
		slog.Error("Error acquiring tab", "err", err)
		os.Exit(1)
	}

	content, err := tab.HTMLContent(url, nil)
	if err != nil {
		slog.Error("Error getting HTML content", "url", url, "err", err)
		os.Exit(1)
	}
	fmt.Println(content)
}

func init() {
	envflags.EnvPrefix = "HEADLESS_"
	browserFlags = cmd.HeadlessBrowserFlags(flags)
	logLevelFlag := envflags.NewLogLevel("LOG_LEVEL", slog.LevelInfo)
	logLevelFlag.AddTo(flags, "log-level", "Log level")
	noHeadlessFlag := envflags.NewBool("NO_HEADLESS", false)
	noHeadlessFlag.AddTo(flags, "H", "Don't run headless mode")
	flags.Usage = usage
	flags.Parse(os.Args[1:])
	slog.SetLogLoggerLevel(logLevelFlag.Get())
	headless = !noHeadlessFlag.Get()
}

func usage() {
	fmt.Println(`Usage: 
	headless [flags] :url
 
  -h	
  	Show this help message`)

	flags.PrintDefaults()
}
