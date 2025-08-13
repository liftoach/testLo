package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testLo/internal/httpserver"
	"testLo/internal/repository"
	"testLo/internal/service"
	"testLo/pkg/server"
	"time"
)

func main() {
	if err := run(); err != nil {
		slog.Error("application error", slog.Any("err", err))
		os.Exit(1)
	}
}

func run() error {
	// Контекст с отменой по сигналу
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	// Канал логирования
	logCh := make(chan string, 100)
	go func() {
		for msg := range logCh {
			log.Printf("[LOG] %s", msg)
		}
	}()

	// Репозиторий и сервис
	taskRepo := repository.NewMemoryTaskRepo()
	svc := service.NewService(taskRepo)

	// Хендлеры
	taskHandler := httpserver.NewTaskHandler(*svc)

	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			taskHandler.GetTasks(w, r)
			logCh <- "GET /tasks called"
		case http.MethodPost:
			taskHandler.CreateTask(w, r)
			logCh <- "POST /tasks called"
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			taskHandler.GetTaskByID(w, r)
			logCh <- "GET /tasks/{id} called"
		} else {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// HTTP сервер
	srv := server.New(":8080", mux)

	// Запуск сервера в отдельной горутине
	errCh := make(chan error, 1)
	go func() {
		slog.Info("server starting", slog.String("addr", ":8080"))
		if err := srv.Run(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	// Ожидание завершения
	select {
	case <-ctx.Done():
		slog.Info("shutdown signal received")
	case err := <-errCh:
		return err
	}

	// Завершение сервера с таймаутом
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}

	close(logCh) // закрываем канал логов
	return nil
}
