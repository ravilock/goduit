package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ravilock/goduit/internal/app/dtos"
	"github.com/ravilock/goduit/internal/app/models"
	"github.com/ravilock/goduit/internal/app/repositories"
	encryptionkeys "github.com/ravilock/goduit/internal/config/encryptionKeys"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

func Register(model *models.User, password string, ctx context.Context) (*models.User, *string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, err
	}
	passwordHashString := string(passwordHash)
	model.PasswordHash = &passwordHashString

	if err = repositories.RegisterUser(model, ctx); err != nil {
		return nil, nil, err
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
		return nil, nil, err
	}

	return model, &tokenString, nil
}
