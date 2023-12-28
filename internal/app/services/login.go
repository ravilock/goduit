package services

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/internal/app/dtos"
	"github.com/ravilock/goduit/internal/app/repositories"
	"github.com/ravilock/goduit/internal/app/transformers"
	encryptionkeys "github.com/ravilock/goduit/internal/config/encryptionKeys"
	"golang.org/x/crypto/bcrypt"
)

func Login(user *dtos.User, ctx context.Context) (*dtos.User, error) {
	model, err := repositories.GetUserByEmail(*user.Email, ctx)
	if err != nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(*model.PasswordHash), []byte(*user.Password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, api.FailedLoginAttempt
		}
		return nil, err
	}

	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, &dtos.TokenClaims{
		Username: *model.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "goduit",
			Subject:   *model.Email,
			ExpiresAt: jwt.NewNumericDate(now.Add(10 * time.Minute)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        uuid.NewString(),
		},
	})

	tokenString, err := token.SignedString(encryptionkeys.PrivateKey)
	if err != nil {
		return nil, err
	}
	user.Token = &tokenString

	return transformers.ModelToUserDto(model, user), nil
}
