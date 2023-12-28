package assemblers

import (
	"errors"

	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/api/responses"
	"github.com/ravilock/goduit/internal/app/dtos"
)

var nilDtoError = errors.New("Dto is nil")
var nilUsernameError = errors.New("Username is nil")

func ProfileResponse(user *dtos.Profile) (*responses.ProfileResponse, error) {
	var profile responses.Profile
	response := new(responses.ProfileResponse)

	if user == nil {
		return nil, api.InternalError(nilDtoError)
	}

	if user.Username == nil {
		return nil, api.InternalError(nilUsernameError)
	}
	profile.Username = *user.Username

	if user.Bio != nil {
		profile.Bio = *user.Bio
	}

	if user.Image != nil {
		profile.Image = *user.Image
	}

	profile.Following = user.Following
	response.Profile = profile
	return response, nil
}
