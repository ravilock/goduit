package services

import (
	"context"

	"github.com/ravilock/goduit/internal/profileManager/models"
)

type getProfileService struct {
	repository UserGetter
}

func (s *getProfileService) GetProfileByEmail(ctx context.Context, email string) (*models.User, error) {
	model, err := s.repository.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (s *getProfileService) GetProfileByUsername(ctx context.Context, username string) (*models.User, error) {
	model, err := s.repository.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return model, nil
}
