package transformers

import (
	"github.com/ravilock/goduit/internal/app/dtos"
	"github.com/ravilock/goduit/internal/app/models"
)

func DtoToModel(user *dtos.User) *models.User {
	return &models.User{
		Username:     user.Username,
		Email:        user.Email,
		Bio:          user.Bio,
		Image:        user.Image,
		PasswordHash: new(string),
	}
}

func ModelToDto(model *models.User, dto *dtos.User) *dtos.User {
	dto.Username = model.Username
	dto.Email = model.Email
	dto.Bio = model.Bio
	dto.Image = model.Image
	return dto
}
