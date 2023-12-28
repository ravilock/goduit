package assemblers

import (
	"github.com/ravilock/goduit/api/responses"
	"github.com/ravilock/goduit/internal/app/dtos"
)

func ProfileResponse(user *dtos.Profile) *responses.Profile {
	response := new(responses.Profile)
	response.Profile.Username = *user.Username
	response.Profile.Following = user.Following
	if user.Bio != nil {
		response.Profile.Bio = *user.Bio
	}
	if user.Image != nil {
		response.Profile.Image = *user.Image
	}
	return response
}
