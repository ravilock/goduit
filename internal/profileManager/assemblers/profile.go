package assemblers

import (
	"errors"

	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/internal/profileManager/models"
	"github.com/ravilock/goduit/internal/profileManager/responses"
)

var (
	errNilModel    = errors.New("model is nil")
	errNilUsername = errors.New("username is nil")
)

func ProfileResponse(user *models.User, isFollowing bool) (*responses.ProfileResponse, error) {
	var profile responses.Profile
	response := new(responses.ProfileResponse)
	if user == nil {
		return nil, api.InternalError(errNilModel)
	}
	if user.Username == nil {
		return nil, api.InternalError(errNilUsername)
	}
	profile.Username = *user.Username
	if user.Bio != nil {
		profile.Bio = *user.Bio
	}
	if user.Image != nil {
		profile.Image = *user.Image
	}
	profile.Following = isFollowing
	response.Profile = profile
	return response, nil
}
