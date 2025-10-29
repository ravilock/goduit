package integrationtests

import (
	"log"
	"os"
	"testing"

	"github.com/ravilock/goduit/internal/config"
	"github.com/ravilock/goduit/internal/mongo"
	"github.com/spf13/viper"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	viper.SetDefault("server.url", "http://localhost:3000")
	if err := config.LoadKeysFromEnv(); err != nil {
		log.Fatal("Failed to load keys from environment variables", err)
	}
	// Connect Mongo DB
	client, err := mongo.ConnectDatabase(viper.GetString("db.url"))
	if err != nil {
		log.Fatalln("Error connecting to database", err)
	}
	ClearDatabase(client)
}
