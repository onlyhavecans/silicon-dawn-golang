package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"onlyhavecans.works/onlyhavecans/silicon-dawn/internal/server"
)

func main() {
	config := &server.Config{
		LogLevel: slog.LevelInfo,
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3200"
	}
	config.Port = port

	cardDir := os.Getenv("CARD_DIR")
	if cardDir == "" {
		cardDir = "data"
	}
	config.CardsDir = cardDir

	srv, err := server.NewServer(config)
	if err != nil {
		slog.Error("failed to build server", slog.Any("error", err))
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := srv.Start(ctx); err != nil {
		slog.Error("failed to start server", slog.Any("error", err))
		os.Exit(1)
	}
}
