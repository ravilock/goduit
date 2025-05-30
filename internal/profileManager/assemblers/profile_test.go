package assemblers

import (
	"testing"

	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/internal/profileManager/models"
)

func TestProfileResponse(t *testing.T) {
	t.Run("Should handle if model is nil", func(t *testing.T) {
		responses, err := ProfileResponse(nil, false)
		if responses != nil {
			t.Errorf("Response should be nil")
		}
		assertError(t, err, api.InternalError(errNilModel))
	})
	t.Run("Should return error if Username is nil", func(t *testing.T) {
		model := assembleUserModel()
		model.Username = nil
		responses, err := ProfileResponse(model, false)
		if responses != nil {
			t.Errorf("Response should be nil")
		}
		assertError(t, err, api.InternalError(errNilUsername))
	})
	t.Run("Should handle if Bio is nil", func(t *testing.T) {
		model := assembleUserModel()
		model.Bio = nil
		responses, err := ProfileResponse(model, false)
		assertNoError(t, err)
		if responses == nil {
			t.Errorf("Response should not be nil")
			return
		}
		if responses.Profile.Bio != "" {
			t.Errorf("Response bio should be blank")
		}
	})
	t.Run("Should handle if Image is nil", func(t *testing.T) {
		model := assembleUserModel()
		model.Image = nil
		responses, err := ProfileResponse(model, false)
		assertNoError(t, err)
		if responses == nil {
			t.Errorf("Response should not be nil")
			return
		}
		if responses.Profile.Image != "" {
			t.Errorf("Response image should be blank")
		}
	})
}

func assembleUserModel() *models.User {
	username := "raylok"
	bio := "This is a bio"
	image := "https://cataas.com/cat"
	return &models.User{
		Username: &username,
		Bio:      &bio,
		Image:    &image,
	}
}

func assertError(t testing.TB, got, want error) {
	t.Helper()
	if got == nil {
		t.Errorf("Expected error but didn't receive one")
		return
	}
	if got.Error() != want.Error() {
		t.Errorf("Got %q want %q", got, want)
		return
	}
}

func assertNoError(t testing.TB, got error) {
	t.Helper()
	if got != nil {
		t.Errorf("Got an error but did'nt want one")
	}
}
