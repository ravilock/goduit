package services

import (
	"context"

	"github.com/ravilock/goduit/internal/app"
	"github.com/ravilock/goduit/internal/app/models"
	"github.com/ravilock/goduit/internal/app/repositories"
	"github.com/ravilock/goduit/internal/identity"
	"golang.org/x/crypto/bcrypt"
)

func Login(model *models.User, password string, ctx context.Context) (*models.User, string, error) {
	model, err := repositories.GetUserByEmail(*model.Email, ctx)
	if err != nil {
		return nil, "", err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(*model.PasswordHash), []byte(password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, "", app.WrongPasswordError.AddContext(err)
		}
		return nil, "", err
	}

	tokenString, err := identity.GenerateToken(model.Username, model.Email)
	if err != nil {
		return nil, "", err
	}

	return model, tokenString, nil
}
