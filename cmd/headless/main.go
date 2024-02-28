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

func init() {
	envflags.EnvPrefix = "HEADLESS_"
	browserFlags = cmd.HeadlessBrowserFlags(flags)
	logLevelFlag := envflags.NewLogLevel("LOG_LEVEL", slog.LevelInfo)
	logLevelFlag.AddTo(flags, "log-level", "Log level")
	flags.Usage = usage
	flags.Parse(os.Args[1:])
	slog.SetLogLoggerLevel(logLevelFlag.Get())
}

func usage() {
	fmt.Println(`Usage: 
	headless [flags] :url
 
  -h	
  	Show this help message`)

	flags.PrintDefaults()
}
