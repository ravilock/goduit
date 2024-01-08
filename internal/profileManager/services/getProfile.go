package services

import (
	"context"

	"github.com/ravilock/goduit/internal/profileManager/models"
)

type UserGetter interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserByID(ctx context.Context, ID string) (*models.User, error)
}

type getProfileService struct {
	repository UserGetter
}

func (s *getProfileService) GetProfileByUsername(ctx context.Context, username string) (*models.User, error) {
	model, err := s.repository.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (s *getProfileService) GetProfileByID(ctx context.Context, ID string) (*models.User, error) {
	model, err := s.repository.GetUserByID(ctx, ID)
	if err != nil {
		return nil, err
	}
	return model, nil
}
