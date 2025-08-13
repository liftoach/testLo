package server

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

}

func run() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	slog.Info("loa")
}

func launchServer(ctx context.Context, srv *http.Server) (err error) {
	var shutdownErr error

	defer func() {
		err = errors.Join(err, shutdownErr)
	}()

	shutdownDone := make(chan struct{})

	go func() {
		<-ctx.Done()
		log.Println("shutting down server...")
		sdCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		shutdownErr = srv.Shutdown(sdCtx)
		log.Println("HTTP server stopped")

		close(shutdownDone)
	}()

	select {
	case <-ctx.Done():
		return nil
	default:
	}
	log.Printf("Listening on %s", srv.Addr)
	if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	<-shutdownDone
	return nil
}
