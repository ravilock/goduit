package services

import (
	"context"

	"github.com/ravilock/goduit/internal/app/dtos"
	"github.com/ravilock/goduit/internal/app/repositories"
	"github.com/ravilock/goduit/internal/app/transformers"
)

func GetProfileByUsername(profileUsername string, ctx context.Context) (*dtos.Profile, error) {
	model, err := repositories.GetUserByUsername(profileUsername, ctx)
	if err != nil {
		return nil, err
	}
	return transformers.ModelToProfileDto(model), nil
}
