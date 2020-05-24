package main

import (
	"log"
	"os"

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

	log.Fatal(server.Start())
}
