package services

import (
	"context"

	"github.com/ravilock/goduit/internal/profileManager/models"
)

type getProfileService struct {
	repository userGetter
}

func (s *getProfileService) GetProfile(ctx context.Context, email string) (*models.User, error) {
	model, err := s.repository.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return model, nil
}
