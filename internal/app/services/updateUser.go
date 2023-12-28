package services

import (
	"context"

	"github.com/ravilock/goduit/internal/app/dtos"
	"github.com/ravilock/goduit/internal/app/repositories"
	"github.com/ravilock/goduit/internal/app/transformers"
	"golang.org/x/crypto/bcrypt"
)

func UpdateUser(user *dtos.User, ctx context.Context) (*dtos.User, error) {
	model := transformers.DtoToModel(user)

	if user.Password != nil {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)
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
	return user, nil
}
