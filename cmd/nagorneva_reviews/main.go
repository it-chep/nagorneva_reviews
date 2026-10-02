package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/nagorneva/nagorneva_reviews/internal"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app, err := internal.New(ctx)
	if err != nil {
		slog.Error("application initialization failed", "error", err)
		os.Exit(1)
	}
	if err := app.Run(ctx); err != nil {
		slog.Error("application stopped with error", "error", err)
		os.Exit(1)
	}
}
