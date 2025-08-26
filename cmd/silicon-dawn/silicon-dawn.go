package main

import (
	"log/slog"
	"os"

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

	server := server.NewServer(config)

	if err := server.Start(); err != nil {
		slog.Error("failed to start server", slog.Any("error", err))
		os.Exit(1)
	}
}
