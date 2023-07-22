package services

import (
	"context"

	"github.com/ravilock/goduit/internal/app/dtos"
	"github.com/ravilock/goduit/internal/app/repositories"
	"github.com/ravilock/goduit/internal/app/transformers"
)

func GetUserByEmail(email string, ctx context.Context) (*dtos.User, error) {
	model, err := repositories.GetUserByEmail(email, ctx)
	if err != nil {
		return nil, err
	}
	return transformers.UserModelToDto(model, new(dtos.User)), nil
}
