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
	"github.com/efixler/headless/browser"
	"github.com/efixler/headless/internal/proxy"
	"github.com/efixler/headless/ua"
	"github.com/efixler/webutil/graceful"
)

var (
	flags         = flag.NewFlagSet("headless-proxy", flag.ExitOnError)
	maxConcurrent *envflags.Value[int]
	userAgent     *envflags.Value[*ua.Arg]
	proxyFlag     = flags.Bool("proxy", false, "Run as a proxy server")
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
		browser.UserAgentIfNotEmpty(userAgent.Get().String()),
	)
	var err error
	if *proxyFlag {
		if server.Handler, err = proxy.HTTPProxy(c); err != nil {
			slog.Error("can't initialize headless proxy", "err", err)
			os.Exit(1)
		}
	} else {
		if server.Handler, err = proxy.Service(c); err != nil {
			slog.Error("can't initialize headless service", "err", err)
			os.Exit(1)
		}
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

	userAgent = envflags.NewText("DEFAULT_USER_AGENT", &ua.Arg{})
	userAgent.AddTo(flags, "default-user-agent", "Default user agent string (omit for browser default, :firefox: for Firefox, :safari: for Safari, or custom string)")
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
