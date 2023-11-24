package identity

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	encryptionkeys "github.com/ravilock/goduit/internal/config/encryptionKeys"
)

var (
	invalidTokenErr        = errors.New("Invalid Token")
	couldNotParseClaimsErr = errors.New("Could Not Parse Claims")
)

type Identyty struct {
	Username string `json:"username,omitempty"`
	jwt.RegisteredClaims
}

func GenerateToken(userEmail, username *string) (string, error) {
	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, &Identyty{
		Username: *username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "goduit",
			Subject:   *userEmail,
			ExpiresAt: jwt.NewNumericDate(now.Add(10 * time.Minute)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        uuid.NewString(),
		},
	})

	return token.SignedString(encryptionkeys.PrivateKey)
}

func FromToken(authorizationHeader string) (*Identyty, error) {
	token, err := jwt.ParseWithClaims(authorizationHeader, &Identyty{}, func(t *jwt.Token) (interface{}, error) {
		key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(os.Getenv("PUBLIC_KEY")))
		if err != nil {
			return nil, err
		}

		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, invalidTokenErr
		}
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, invalidTokenErr
	}
	claims, ok := token.Claims.(*Identyty)
	if !ok {
		return nil, couldNotParseClaimsErr
	}
	return claims, nil
}
