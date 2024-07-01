package services

import (
	"context"
	"errors"

	"github.com/ravilock/goduit/internal/identity"
	"github.com/ravilock/goduit/internal/password"
	"github.com/ravilock/goduit/internal/profileManager/models"
)

type UserRegister interface {
	RegisterUser(ctx context.Context, user *models.User) (*models.User, error)
}

type registerProfileService struct {
	repository UserRegister
	hasher     password.Hasher
}

var (
	ErrFailedToGeneratePasswordHash = errors.New("failed to generate password hash")
	ErrFailedToRegisterUser         = errors.New("failed to register user")
)

func (s *registerProfileService) Register(ctx context.Context, model *models.User, password string) (string, error) {
	passwordHash, err := s.hasher.Hash(password)
	if err != nil {
		return "", errors.Join(ErrFailedToGeneratePasswordHash, err)
	}
	passwordHashString := string(passwordHash)
	model.PasswordHash = &passwordHashString

	model, err = s.repository.RegisterUser(ctx, model)
	if err != nil {
		return "", errors.Join(ErrFailedToRegisterUser, err)
	}

	tokenString, err := identity.GenerateToken(*model.Email, *model.Username, model.ID.Hex())
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
