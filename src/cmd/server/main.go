package server

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {

}

func run() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	slog.Info("loa")
}
