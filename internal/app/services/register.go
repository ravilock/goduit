package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ravilock/goduit/internal/app/dtos"
	"github.com/ravilock/goduit/internal/app/repositories"
	"github.com/ravilock/goduit/internal/app/transformers"
	encryptionkeys "github.com/ravilock/goduit/internal/config/encryptionKeys"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

func Register(user *dtos.User, ctx context.Context) (*dtos.User, error) {
	model := transformers.DtoToModel(user)

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	passwordHashString := string(passwordHash)
	model.PasswordHash = &passwordHashString

	if err = repositories.RegisterUser(model, ctx); err != nil {
		return nil, err
	}

	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, &jwt.RegisteredClaims{
		Issuer:    "goduit",
		Subject:   *model.Email,
		ExpiresAt: jwt.NewNumericDate(now.Add(10 * time.Minute)),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        uuid.NewString(),
	})

	tokenString, err := token.SignedString(encryptionkeys.PrivateKey)
	if err != nil {
		return nil, err
	}
	user.Token = &tokenString

	return user, nil
}
