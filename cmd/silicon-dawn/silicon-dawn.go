package main

import (
	"os"

	"github.com/rs/zerolog/log"

	"onlyhavecans.works/amy/silicondawn"
)

func main() {
	config := &silicondawn.Config{}

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

	server := silicondawn.NewServer(config)

	if err := server.Start(); err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}
