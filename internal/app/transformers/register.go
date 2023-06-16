package transformers

import (
	"github.com/ravilock/goduit/internal/app/dtos"
	"github.com/ravilock/goduit/internal/app/models"
)

func DtoToModel(user *dtos.User) *models.User {
	return &models.User{
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: new(string),
		Bio:          new(string),
		Image:        new(string),
	}
}
