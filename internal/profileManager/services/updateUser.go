package services

import (
	"context"

	"github.com/ravilock/goduit/internal/identity"
	"github.com/ravilock/goduit/internal/profileManager/models"
	"golang.org/x/crypto/bcrypt"
)

type profileUpdater interface {
	UpdateProfile(ctx context.Context, subjectEmail, clientUsername string, user *models.User) error
}

type updateProfileService struct {
	repository profileUpdater
}

func (s *updateProfileService) UpdateProfile(ctx context.Context, subjectEmail, clientUsername, password string, model *models.User) (string, error) {
	if shouldGenerateNewPasswordHash(password) {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return "", err
		}
		passwordHashString := string(passwordHash)
		model.PasswordHash = &passwordHashString
	}

	err := s.repository.UpdateProfile(ctx, subjectEmail, clientUsername, model)
	if err != nil {
		return "", err
	}

	var token string
	if shouldGenerateNewToken(subjectEmail, clientUsername, model) {
		token, err = identity.GenerateToken(*model.Email, *model.Username, model.ID.Hex())
		if err != nil {
			return "", err
		}
	}
	return token, nil
}

func shouldGenerateNewPasswordHash(password string) bool {
	return password != ""
}

func shouldGenerateNewToken(subjectEmail, clientUsername string, model *models.User) bool {
	return subjectEmail != *model.Email || clientUsername != *model.Username
}
