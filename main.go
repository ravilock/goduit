package main

import (
	"log"
	"os"

	"github.com/ravilock/goduit/internal/api"
	"github.com/ravilock/goduit/internal/config"
	"github.com/spf13/viper"
)

func main() {
	privateKeyFile, err := os.Open(viper.GetString("private.key.location"))
	if err != nil {
		log.Fatal("Failed to open private key file", err)
	}

	if err := config.LoadPrivateKey(privateKeyFile); err != nil {
		log.Fatal("Failed to load private key file content", err)
	}

	if err := privateKeyFile.Close(); err != nil {
		log.Fatal("Failed to close private key file", err)
	}

	publicKeyFile, err := os.Open(viper.GetString("public.key.location"))
	if err != nil {
		log.Fatal("Failed to open public key file", err)
	}

	if err := config.LoadPublicKey(publicKeyFile); err != nil {
		log.Fatal("Failed to load public key file content", err)
	}

	if err := publicKeyFile.Close(); err != nil {
		log.Fatal("Failed to close publicKeyFile key file", err)
	}

	server, err := api.NewServer()
	if err != nil {
		log.Fatalln("Failed to start server", err)
	}
	server.Start()
}
