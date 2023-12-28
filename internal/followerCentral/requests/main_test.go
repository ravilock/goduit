package requests

import (
	"log"
	"math/rand"
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

func randomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
