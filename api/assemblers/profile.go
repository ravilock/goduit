package assemblers

import (
	"github.com/ravilock/goduit/api/responses"
	"github.com/ravilock/goduit/internal/app/dtos"
)

func ProfileResponse(user *dtos.Profile) *responses.ProfileResponse {
	var profile responses.Profile
	response := new(responses.ProfileResponse)
	profile.Username = *user.Username
	profile.Following = user.Following
	if user.Bio != nil {
		profile.Bio = *user.Bio
	}
	if user.Image != nil {
		profile.Image = *user.Image
	}
	response.Profile = profile
	return response
}
