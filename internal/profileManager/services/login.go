package services

import (
	"context"

	"github.com/ravilock/goduit/internal/app"
	"github.com/ravilock/goduit/internal/identity"
	"github.com/ravilock/goduit/internal/profileManager/models"
	"golang.org/x/crypto/bcrypt"
)

type LogUserService struct {
	repository UserGetter
}

func NewLogUserService(repository UserGetter) *LogUserService {
	return &LogUserService{
		repository: repository,
	}
}

func (s *LogUserService) Login(ctx context.Context, email, password string) (*models.User, string, error) {
	model, err := s.repository.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, "", err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(*model.PasswordHash), []byte(password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, "", app.WrongPasswordError.AddContext(err)
		}
		return nil, "", err
	}

	tokenString, err := identity.GenerateToken(*model.Email, *model.Username, model.ID.Hex())
	if err != nil {
		return nil, "", err
	}

	return model, tokenString, nil
}
