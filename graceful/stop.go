package graceful

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// After a SIGINT, SIGTERM, or os.Interrupt, shut down s, wait a bit for
// requests to clear, then call the passed cancel function.
// This will let the requests finish before shutting down any other open
// resources with the passed CancelFunc.
// This function will block until server shutdown is complete
func WaitForShutdown(s *http.Server, cf context.CancelFunc) {
	waitForKill()
	<-shutdownServer(s, cf)
}

func waitForKill() {
	kill := make(chan os.Signal, 1)
	signal.Notify(kill, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-kill
}

// Shutdown the server and then propagate the shutdown to a cancel function.

// Caller should block on the returned channel.
func shutdownServer(s *http.Server, cf context.CancelFunc) chan bool {
	slog.Info("server shutting down")
	wchan := make(chan bool)
	// a large request set can take a while to finish,
	// so we give the server a couple minutes to finish if it needs to
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	context.AfterFunc(ctx, func() {
		cf()
		// without a little bit of sleep here sometimes final log messages
		// don't get flushed, even with the file sync above
		time.Sleep(100 * time.Millisecond)
		close(wchan)
	})
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		slog.Error("server shutdown failed", "error", err)
	}
	slog.Info("server stopped")
	return wchan
}
