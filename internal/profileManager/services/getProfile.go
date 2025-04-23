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

type GetProfileService struct {
	repository UserGetter
}

func NewGetProfileService(repository UserGetter) *GetProfileService {
	return &GetProfileService{
		repository: repository,
	}
}

func (s *GetProfileService) GetProfileByUsername(ctx context.Context, username string) (*models.User, error) {
	model, err := s.repository.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (s *GetProfileService) GetProfileByID(ctx context.Context, ID string) (*models.User, error) {
	model, err := s.repository.GetUserByID(ctx, ID)
	if err != nil {
		return nil, err
	}
	return model, nil
}
