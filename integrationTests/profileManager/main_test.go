package profilemanager

import (
	"log"
	"os"
	"testing"

	"github.com/ravilock/goduit/internal/config"
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
}
