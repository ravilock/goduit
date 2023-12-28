package services

import (
	"context"

	"github.com/ravilock/goduit/internal/app/dtos"
	"github.com/ravilock/goduit/internal/app/repositories"
	"github.com/ravilock/goduit/internal/app/transformers"
)

func GetProfileByUsername(username, loggedUsername string, ctx context.Context) (*dtos.Profile, error) {
	model, err := repositories.GetUserByUsername(username, ctx)
	if err != nil {
		return nil, err
	}
	following := IsFollowedBy(username, loggedUsername, ctx)
	return transformers.ModelToProfileDto(model, following), nil
}
