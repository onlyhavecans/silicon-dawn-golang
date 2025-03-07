package main

import (
	"os"

	"github.com/onlyhavecans/silicondawn/internal/server"
	"github.com/rs/zerolog/log"
)

func main() {
	config := &server.Config{}

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
		log.Fatal().Err(err).Msg("failed to start server")
	}
}
