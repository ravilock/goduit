package services

import (
	"github.com/ravilock/goduit/internal/app/models"
	"github.com/ravilock/goduit/internal/app/repositories"
	"github.com/ravilock/goduit/internal/identity"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

func Register(model *models.User, password string, ctx context.Context) (*models.User, string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}
	passwordHashString := string(passwordHash)
	model.PasswordHash = &passwordHashString

	if err = repositories.RegisterUser(model, ctx); err != nil {
		return nil, "", err
	}

	tokenString, err := identity.GenerateToken(model.Username, model.Email)
	if err != nil {
		return nil, "", err
	}

	return model, tokenString, nil
}
