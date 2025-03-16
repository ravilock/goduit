package assemblers

import (
	"github.com/ravilock/goduit/internal/profileManager/models"
	"github.com/ravilock/goduit/internal/profileManager/responses"
)

func UserResponse(user *models.User) *responses.User {
	response := new(responses.User)
	if user.Username != nil {
		response.User.Username = *user.Username
	}
	if user.Email != nil {
		response.User.Email = *user.Email
	}
	if user.Bio != nil {
		response.User.Bio = *user.Bio
	}
	if user.Image != nil {
		response.User.Image = *user.Image
	}
	return response
}
