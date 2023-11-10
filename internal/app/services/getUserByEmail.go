package services

import (
	"context"

	"github.com/ravilock/goduit/internal/app/models"
	"github.com/ravilock/goduit/internal/app/repositories"
)

func GetUserByEmail(email string, ctx context.Context) (*models.User, error) {
	model, err := repositories.GetUserByEmail(email, ctx)
	if err != nil {
		return nil, err
	}
	return model, nil
}
