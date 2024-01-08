package services

import (
	"context"

	"github.com/ravilock/goduit/internal/identity"
	"github.com/ravilock/goduit/internal/profileManager/models"
	"golang.org/x/crypto/bcrypt"
)

type profileRegister interface {
	RegisterUser(ctx context.Context, user *models.User) (*models.User, error)
}

type registerProfileService struct {
	repository profileRegister
}

func (s *registerProfileService) Register(ctx context.Context, model *models.User, password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	passwordHashString := string(passwordHash)
	model.PasswordHash = &passwordHashString

	model, err = s.repository.RegisterUser(ctx, model)
	if err != nil {
		return "", err
	}

	tokenString, err := identity.GenerateToken(*model.Email, *model.Username, model.ID.Hex())
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
