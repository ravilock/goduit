package requests

import (
	"testing"

	"github.com/ravilock/goduit/api"
	"github.com/stretchr/testify/require"
)

func TestUpdateProfile(t *testing.T) {
	t.Run("Valid request should not return errors", func(t *testing.T) {
		request := generateUpdateProfileRequest()
		err := request.Validate()
		require.NoError(t, err)
	})

	t.Run("Email is required", func(t *testing.T) {
		request := generateUpdateProfileRequest()
		request.User.Email = ""
		err := request.Validate()
		require.ErrorContains(t, err, api.RequiredFieldError("Email").Error())
	})

	t.Run("Username is required", func(t *testing.T) {
		request := generateUpdateProfileRequest()
		request.User.Username = ""
		err := request.Validate()
		require.ErrorContains(t, err, api.RequiredFieldError("Username").Error())
	})

	t.Run("Bio is required", func(t *testing.T) {
		request := generateUpdateProfileRequest()
		request.User.Bio = ""
		err := request.Validate()
		require.ErrorContains(t, err, api.RequiredFieldError("Bio").Error())
	})

	t.Run("Image is required", func(t *testing.T) {
		request := generateUpdateProfileRequest()
		request.User.Image = ""
		err := request.Validate()
		require.ErrorContains(t, err, api.RequiredFieldError("Image").Error())
	})

	t.Run("Email should not be blank", func(t *testing.T) {
		request := generateUpdateProfileRequest()
		request.User.Email = " "
		err := request.Validate()
		require.ErrorContains(t, err, api.RequiredFieldError("Email").Error())
	})

	t.Run("Password should not be blank", func(t *testing.T) {
		request := generateUpdateProfileRequest()
		request.User.Password = " "
		err := request.Validate()
		require.ErrorContains(t, err, api.RequiredFieldError("Password").Error())
	})

	t.Run("Username should not be blank", func(t *testing.T) {
		request := generateUpdateProfileRequest()
		request.User.Username = " "
		err := request.Validate()
		require.ErrorContains(t, err, api.RequiredFieldError("Username").Error())
	})

	t.Run("Bio should not be blank", func(t *testing.T) {
		request := generateUpdateProfileRequest()
		request.User.Bio = " "
		err := request.Validate()
		require.ErrorContains(t, err, api.RequiredFieldError("Bio").Error())
	})

	t.Run("Image should not be blank", func(t *testing.T) {
		request := generateUpdateProfileRequest()
		request.User.Image = " "
		err := request.Validate()
		require.ErrorContains(t, err, api.RequiredFieldError("Image").Error())
	})

	t.Run("Password should contain at least 8 chars", func(t *testing.T) {
		request := generateUpdateProfileRequest()
		request.User.Password = "pass"
		err := request.Validate()
		require.ErrorContains(t, err, api.InvalidFieldLength("Password", "min", "8").Error())
	})

	t.Run("Username should contain at least 5 chars", func(t *testing.T) {
		request := generateUpdateProfileRequest()
		request.User.Username = "user"
		err := request.Validate()
		require.ErrorContains(t, err, api.InvalidFieldLength("Username", "min", "5").Error())
	})

	t.Run("Email should contain at most 256 chars", func(t *testing.T) {
		request := generateUpdateProfileRequest()
		request.User.Email = randomString(256) + "@hotmail.com"
		err := request.Validate()
		require.ErrorContains(t, err, api.InvalidFieldLength("Email", "max", "256").Error())
	})

	t.Run("Password should contain at most 72 chars", func(t *testing.T) {
		request := generateUpdateProfileRequest()
		request.User.Password = randomString(73)
		err := request.Validate()
		require.ErrorContains(t, err, api.InvalidFieldLength("Password", "max", "72").Error())
	})

	t.Run("Username should contain at most 255 chars", func(t *testing.T) {
		request := generateUpdateProfileRequest()
		request.User.Username = randomString(256)
		err := request.Validate()
		require.ErrorContains(t, err, api.InvalidFieldLength("Username", "max", "255").Error())
	})

	t.Run("Email should be a valid email", func(t *testing.T) {
		request := generateUpdateProfileRequest()
		request.User.Email = "email@hotmail."
		err := request.Validate()
		require.ErrorContains(t, err, api.InvalidFieldError("Email", request.User.Email).Error())
	})
}

func generateUpdateProfileRequest() *UpdateProfileRequest {
	updateProfile := new(UpdateProfileRequest)
	updateProfile.User.Username = "test-username"
	updateProfile.User.Email = "test-email@email.com"
	updateProfile.User.Password = "test-password"
	updateProfile.User.Bio = "Test Bio"
	updateProfile.User.Image = "https://update.profile.test.bio.com/image"
	return updateProfile
}
