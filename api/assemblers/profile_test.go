package assemblers

import (
	"testing"

	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/internal/app/dtos"
)

func TestProfileResponse(t *testing.T) {
	t.Run("Should handle if dto is nil", func(t *testing.T) {
		responses, err := ProfileResponse(nil)
		if responses != nil {
			t.Errorf("Response should be nil")
		}
		assertError(t, err, api.InternalError(nilDtoError))

	})
	t.Run("Should return error if Username is nil", func(t *testing.T) {
		dto := assebleProfileDto()
		dto.Username = nil
		responses, err := ProfileResponse(dto)
		if responses != nil {
			t.Errorf("Response should be nil")
		}
		assertError(t, err, api.InternalError(nilUsernameError))
	})
	t.Run("Should handle if Bio is nil", func(t *testing.T) {
		dto := assebleProfileDto()
		dto.Bio = nil
		responses, err := ProfileResponse(dto)
		assertNoError(t, err)
		if responses == nil {
			t.Errorf("Response should not be nil")
		}
		if responses.Profile.Bio != "" {
			t.Errorf("Response bio should be blank")
		}
	})
	t.Run("Should handle if Image is nil", func(t *testing.T) {
		dto := assebleProfileDto()
		dto.Image = nil
		responses, err := ProfileResponse(dto)
		assertNoError(t, err)
		if responses == nil {
			t.Errorf("Response should not be nil")
		}
		if responses.Profile.Image != "" {
			t.Errorf("Response image should be blank")
		}
	})

}

func assebleProfileDto() *dtos.Profile {
	username := "raylok"
	bio := "This is a bio"
	image := "https://cataas.com/cat"
	return &dtos.Profile{
		Username:  &username,
		Bio:       &bio,
		Image:     &image,
		Following: true,
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
