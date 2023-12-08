package requests

import (
	"testing"

	"github.com/ravilock/goduit/api"
	"github.com/stretchr/testify/assert"
)

func TestGetProfile(t *testing.T) {
	t.Run("Valid request should not return errors", func(t *testing.T) {
		request := generateGetProfileRequest()
		err := request.Validate()
		assert.NoError(t, err)
	})

	t.Run("Username is required", func(t *testing.T) {
		request := generateGetProfileRequest()
		request.Username = ""
		err := request.Validate()
		assert.ErrorContains(t, err, api.RequiredFieldError("Username").Error())
	})

	t.Run("Username should not be blank", func(t *testing.T) {
		request := generateGetProfileRequest()
		request.Username = " "
		err := request.Validate()
		assert.ErrorContains(t, err, api.RequiredFieldError("Username").Error())
	})

	t.Run("Username should contain at least 5 chars", func(t *testing.T) {
		request := generateGetProfileRequest()
		request.Username = "user"
		err := request.Validate()
		assert.ErrorContains(t, err, api.InvalidFieldLength("Username", "min", "5").Error())
	})

	t.Run("Username should contain at most 72 chars", func(t *testing.T) {
		request := generateGetProfileRequest()
		request.Username = randomString(256)
		err := request.Validate()
		assert.ErrorContains(t, err, api.InvalidFieldLength("Username", "max", "255").Error())
	})
}

func generateGetProfileRequest() *GetProfile {
	getProfile := new(GetProfile)
	getProfile.Username = "test.username"
	return getProfile
}
