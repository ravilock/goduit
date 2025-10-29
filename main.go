package main

import (
	"log"

	"github.com/ravilock/goduit/internal/api"
	"github.com/ravilock/goduit/internal/config"
)

func main() {
	if err := config.LoadKeysFromEnv(); err != nil {
		log.Fatal("Failed to load keys from environment variables", err)
	}

	server, err := api.NewServer()
	if err != nil {
		log.Fatalln("Failed to start server", err)
	}
	server.Start()
}
