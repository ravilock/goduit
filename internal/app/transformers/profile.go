package transformers

import (
	"github.com/ravilock/goduit/internal/app/dtos"
	"github.com/ravilock/goduit/internal/app/models"
)

func ModelToProfileDto(model *models.User) *dtos.Profile {
	return &dtos.Profile{
		Username: model.Username,
		Bio:      model.Bio,
		Image:    model.Image,
	}
}
