package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/aachex/taskrunner/server"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT)

	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))

	srv := server.New(logger)
	go srv.Start()

	// shutdown
	<-interrupt
	logger.Info("shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Info("shutdown error", "message", err.Error())
	}
}
