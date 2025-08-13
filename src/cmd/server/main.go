package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testLo/internal/httpserver"
	"testLo/internal/repository"
	"testLo/internal/service"
	"testLo/pkg/logger/slogpretty"
	"testLo/pkg/server"
	"time"
)

const serverPort = ":8080" // <--- вынес порт в константу

func main() {
	if err := run(); err != nil {
		slog.Error("application error", slog.Any("err", err))
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	log := setupLogger()

	logCh := make(chan string, 100)
	go func() {
		for msg := range logCh {
			log.Info(fmt.Sprintf("[LOG] %v", msg))
		}
	}()

	taskRepo := repository.NewMemoryTaskRepo(logCh)
	svc := service.NewService(taskRepo, logCh)

	taskHandler := httpserver.NewTaskHandler(*svc, logCh)

	mux := http.NewServeMux()
	taskHandler.HandleRoutes(mux, logCh)

	srv := server.New(serverPort, mux) // используем константу

	errCh := make(chan error, 1)
	go func() {
		log.Info("server starting", slog.String("addr", serverPort)) // и здесь
		if err := srv.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		slog.Info("shutdown signal received")
	case err := <-errCh:
		return err
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}

	close(logCh)
	return nil
}

func setupLogger() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelInfo,
		},
	}
	handler := opts.NewPrettyHandler(os.Stdout)
	return slog.New(handler)
}
