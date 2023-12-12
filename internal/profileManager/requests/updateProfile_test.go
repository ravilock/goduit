package requests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateProfile(t *testing.T) {
	t.Run("Valid request should not return errors", func(t *testing.T) {
		request := generateUpdateProfileRequest()
		err := request.Validate()
		assert.NoError(t, err)
	})
	// TODO: add checks for theses values as they are optional and can be empty
}

func generateUpdateProfileRequest() *UpdateProfile {
	updateProfile := new(UpdateProfile)
	updateProfile.User.Username = "test-username"
	updateProfile.User.Email = "test-email@email.com"
	updateProfile.User.Password = "test-password"
	updateProfile.User.Bio = "Test Bio"
	updateProfile.User.Image = "https://update.profile.test.bio.com/image"
	return updateProfile
}
