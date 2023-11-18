package services

import (
	"context"

	"github.com/ravilock/goduit/internal/profileManager/models"
	"golang.org/x/crypto/bcrypt"
)

type profileUpdater interface {
	UpdateProfile(ctx context.Context, user *models.User) (*models.User, error)
}

type updateProfileService struct {
	repository profileUpdater
}

func (s *updateProfileService) UpdateProfile(ctx context.Context, model *models.User, password string) (*models.User, error) {
	if password != "" {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		passwordHashString := string(passwordHash)
		model.PasswordHash = &passwordHashString
	}

	model, err := s.repository.UpdateProfile(ctx, model)
	if err != nil {
		return nil, err
	}
	return model, nil
}
