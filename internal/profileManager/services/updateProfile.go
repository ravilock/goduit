package services

import (
	"context"
	"fmt"

	"github.com/ravilock/goduit/internal/identity"
	"github.com/ravilock/goduit/internal/profileManager/models"
	"golang.org/x/crypto/bcrypt"
)

type profileUpdater interface {
	UpdateProfile(ctx context.Context, subjectEmail, clientUsername string, user *models.User) (*models.User, error)
}

type updateProfileService struct {
	repository profileUpdater
}

func (s *updateProfileService) UpdateProfile(ctx context.Context, subjectEmail, clientUsername, password string, model *models.User) (*models.User, string, error) {
	if shouldGenerateNewPasswordHash(password) {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		fmt.Println(passwordHash, password)
		if err != nil {
			return nil, "", err
		}
		passwordHashString := string(passwordHash)
		model.PasswordHash = &passwordHashString
	}

	model, err := s.repository.UpdateProfile(ctx, subjectEmail, clientUsername, model)
	if err != nil {
		return nil, "", err
	}

	if shouldGenerateNewToken(subjectEmail, clientUsername, model) {
		tokenString, err := identity.GenerateToken(*model.Email, *model.Username)
		if err != nil {
			return nil, "", err
		}
		return model, tokenString, nil
	}
	return model, "", nil
}

func shouldGenerateNewPasswordHash(password string) bool {
	return password != ""
}

func shouldGenerateNewToken(subjectEmail, clientUsername string, model *models.User) bool {
	return subjectEmail != *model.Email || clientUsername != *model.Username
}
