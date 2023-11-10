package services

import (
	"context"

	"github.com/ravilock/goduit/internal/app/models"
	"github.com/ravilock/goduit/internal/app/repositories"
)

func GetProfileByUsername(profileUsername string, ctx context.Context) (*models.User, error) {
	model, err := repositories.GetUserByUsername(profileUsername, ctx)
	if err != nil {
		return nil, err
	}
	return model, nil
}
