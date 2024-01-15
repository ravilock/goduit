package requests

import (
	"testing"

	"github.com/ravilock/goduit/api"
	"github.com/stretchr/testify/require"
)

func TestFollower(t *testing.T) {
	t.Run("Valid request should not return errors", func(t *testing.T) {
		request := generateFollowerRequest()
		err := request.Validate()
		require.NoError(t, err)
	})
	t.Run("Username is required", func(t *testing.T) {
		request := generateFollowerRequest()
		request.Username = ""
		err := request.Validate()
		require.ErrorContains(t, err, api.RequiredFieldError("Username").Error())
	})
	t.Run("Username should not be blank", func(t *testing.T) {
		request := generateFollowerRequest()
		request.Username = " "
		err := request.Validate()
		require.ErrorContains(t, err, api.RequiredFieldError("Username").Error())
	})
	t.Run("Username should contain at least 5 chars", func(t *testing.T) {
		request := generateFollowerRequest()
		request.Username = "1234"
		err := request.Validate()
		require.ErrorContains(t, err, api.InvalidFieldLimit("Username", "min", "5").Error())
	})
	t.Run("Username should contain at most 255 chars", func(t *testing.T) {
		request := generateFollowerRequest()
		request.Username = randomString(256)
		err := request.Validate()
		require.ErrorContains(t, err, api.InvalidFieldLimit("Username", "max", "255").Error())
	})
}

func generateFollowerRequest() *FollowerRequest {
	return &FollowerRequest{
		Username: "test-username",
	}
}
