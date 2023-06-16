package services

import (
	"github.com/ravilock/goduit/internal/app/dtos"
	"github.com/ravilock/goduit/internal/app/repositories"
	"github.com/ravilock/goduit/internal/app/transformers"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

func Register(user *dtos.User, ctx context.Context) (*dtos.User, error) {
	model := transformers.DtoToModel(user)

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	*model.PasswordHash = string(passwordHash)

	if err = repositories.RegisterUser(model, ctx); err != nil {
		return nil, err
	}

	return user, nil
}
