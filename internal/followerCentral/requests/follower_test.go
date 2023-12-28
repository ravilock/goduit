package requests

import (
	"math/rand"
	"testing"

	"github.com/ravilock/goduit/api"
	"github.com/stretchr/testify/assert"
)

func TestFollower(t *testing.T) {
	t.Run("Valid request should not return errors", func(t *testing.T) {
		request := generateFollowerRequest()
		err := request.Validate()
		assert.NoError(t, err)
	})
	t.Run("Username is required", func(t *testing.T) {
		request := generateFollowerRequest()
		request.Username = ""
		err := request.Validate()
		assert.ErrorContains(t, err, api.RequiredFieldError("Username").Error())
	})
	t.Run("Username should not be blank", func(t *testing.T) {
		request := generateFollowerRequest()
		request.Username = " "
		err := request.Validate()
		assert.ErrorContains(t, err, api.RequiredFieldError("Username").Error())
	})
	t.Run("Username should contain at least 5 chars", func(t *testing.T) {
		request := generateFollowerRequest()
		request.Username = "1234"
		err := request.Validate()
		assert.ErrorContains(t, err, api.InvalidFieldLength("Username", "min", "5").Error())
	})
	t.Run("Username should contain at most 255 chars", func(t *testing.T) {
		request := generateFollowerRequest()
		request.Username = randomString(256)
		err := request.Validate()
		assert.ErrorContains(t, err, api.InvalidFieldLength("Username", "max", "255").Error())
	})
}

func generateFollowerRequest() *Follower {
	return &Follower{
		Username: "test-username",
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
