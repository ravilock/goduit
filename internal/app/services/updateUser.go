package services

import (
	"context"

	"github.com/ravilock/goduit/internal/app/models"
	"github.com/ravilock/goduit/internal/app/repositories"
	"golang.org/x/crypto/bcrypt"
)

func UpdateUser(model *models.User, password string, ctx context.Context) (*models.User, error) {
	if password != "" {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		passwordHashString := string(passwordHash)
		model.PasswordHash = &passwordHashString
	}

	model, err := repositories.UpdateUser(model, ctx)
	if err != nil {
		return nil, err
	}
	return model, nil
}
