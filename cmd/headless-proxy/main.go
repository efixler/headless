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
	"github.com/efixler/headless/serverutil"
)

var (
	flags     = flag.NewFlagSet("headless-proxy", flag.ExitOnError)
	server    = &http.Server{}
	logWriter io.Writer
)

func main() {
	slog.Info("Starting headless-proxy server", "addr", server.Addr)
	server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	})

	_, cancel := context.WithCancel(context.Background())

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("headless-proxy error, shutting down", "error", err)
		}
	}()

	serverutil.WaitForShutdown(server, cancel)
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
	readTimeout.AddTo(flags, "read-timeout", "Read timeout")
	writeTimeout := envflags.NewDuration("WRITE_TIMEOUT", 30*time.Second)
	writeTimeout.AddTo(flags, "write-timeout", "Write timeout")
	idleTimeout := envflags.NewDuration("IDLE_TIMEOUT", 120*time.Second)
	idleTimeout.AddTo(flags, "idle-timeout", "Idle timeout")
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
