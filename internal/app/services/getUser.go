package services

import (
	"context"

	"github.com/ravilock/goduit/internal/app/dtos"
	"github.com/ravilock/goduit/internal/app/repositories"
	"github.com/ravilock/goduit/internal/app/transformers"
)

func GetUser(user *dtos.User, ctx context.Context) (*dtos.User, error) {
	model, err := repositories.GetUserByEmail(*user.Email, ctx)
	if err != nil {
		return nil, err
	}
	return transformers.ModelToDto(model, user), nil
}
