package serverutil

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func WaitForKill() {
	kill := make(chan os.Signal, 1)
	signal.Notify(kill, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-kill
}

// Shutdown the server and then progate the shutdown to the mux
// This will let the requests finish before shutting down the db
// cf is the cancel function for the mux context, or, generically
// speaking, a cancel function to queue up after the server is done
// Caller should block on the returned channel.
func ShutdownServer(s *http.Server, cf context.CancelFunc) chan bool {
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
