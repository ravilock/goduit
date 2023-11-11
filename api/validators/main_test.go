package validators

import (
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	if err := InitValidator(); err != nil {
		log.Fatalln("Failed to load validator", err)
	}
}
