package requests

import (
	"log"
	"os"
	"testing"

	"github.com/ravilock/goduit/api/validators"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	if err := validators.InitValidator(); err != nil {
		log.Fatalln("Failed to load validator", err)
	}
}
