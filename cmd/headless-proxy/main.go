package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/efixler/envflags"
	"github.com/efixler/headless/graceful"
	"github.com/efixler/headless/internal/browser"
	"github.com/efixler/headless/internal/proxy"
)

var (
	flags         = flag.NewFlagSet("headless-proxy", flag.ExitOnError)
	maxConcurrent *envflags.Value[int]
	userAgent     *envflags.Value[string]
	server        = &http.Server{}
	logWriter     io.Writer
)

func main() {
	slog.Info("Starting headless-proxy server", "addr", server.Addr)
	ctx, cancel := context.WithCancel(context.Background())

	c := browser.NewChrome(
		ctx,
		browser.Headless(true),
		browser.MaxTabs(maxConcurrent.Get()),
		browser.UserAgentIfNotEmpty(userAgent.Get()),
	)
	var err error
	if server.Handler, err = proxy.New(c); err != nil {
		slog.Error("can't initialize proxy", "err", err)
		os.Exit(1)
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("headless-proxy error, shutting down", "error", err)
		}
	}()

	graceful.WaitForShutdown(server, cancel)
	if logFile, ok := (logWriter).(*os.File); ok {
		logFile.Sync()
	}
}

func init() {
	logWriter = os.Stderr
	envflags.EnvPrefix = "HEADLESS_PROXY_"
	flags.Usage = usage
	port := envflags.NewInt("PORT", 8008)
	port.AddTo(flags, "port", "Port to listen on")
	readTimeout := envflags.NewDuration("READ_TIMEOUT", 5*time.Second)
	readTimeout.AddTo(flags, "inbound-read-timeout", "Inbound connection read timeout")
	writeTimeout := envflags.NewDuration("WRITE_TIMEOUT", 30*time.Second)
	writeTimeout.AddTo(flags, "inbound-write-timeout", "Inbound connection write timeout")
	idleTimeout := envflags.NewDuration("IDLE_TIMEOUT", 120*time.Second)
	idleTimeout.AddTo(flags, "inbound-idle-timeout", "Inbound connection keepalive idle timeout")
	maxConcurrent = envflags.NewInt("MAX_CONCURRENT", 6)
	maxConcurrent.AddTo(flags, "max-concurrent", "Maximum concurrent connections")
	userAgent = envflags.NewString("DEFAULT_USER_AGENT", "")
	userAgent.AddTo(flags, "default-user-agent", "Default user agent string (empty for browser default)")
	logLevel := envflags.NewLogLevel("LOG_LEVEL", slog.LevelInfo)
	logLevel.AddTo(flags, "log-level", "Set the log level [debug|error|info|warn]")
	flags.Parse(os.Args[1:])
	server.Addr = fmt.Sprintf(":%d", port.Get())
	server.ReadTimeout = readTimeout.Get()
	server.WriteTimeout = writeTimeout.Get()
	server.IdleTimeout = idleTimeout.Get()
	logger := slog.New(slog.NewTextHandler(
		logWriter,
		&slog.HandlerOptions{
			Level: logLevel.Get(),
		},
	))
	slog.SetDefault(logger)
}

func usage() {
	fmt.Println(`Usage: 
	headless-proxy [flags] :url
 
  -h	
  	Show this help message`)

	flags.PrintDefaults()
}
