package identity

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
	encryptionkeys "github.com/ravilock/goduit/internal/config/encryptionKeys"
)

var (
	invalidTokenErr        = errors.New("Invalid Token")
	couldNotParseClaimsErr = errors.New("Could Not Parse Claims")
)

type Identity struct {
	UserEmail string `json:"userId,omitempty"`
	Username  string `json:"username,omitempty"`
	jwt.RegisteredClaims
}

type IdentityHeaders struct {
	Subject      string `header:"Goduit-Subject"`
	ClientUsername string `header:"Goduit-Client-Username"`
	ClientEmail    string `header:"Goduit-Client-Email"`
}

func CreateAuthMiddleware(requiredAuthentication bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authorizationHeader := c.Request().Header.Get("Authorization")
			if !requiredAuthentication && authorizationHeader == "" {
				return next(c)
			}
			identity, err := FromToken(authorizationHeader)
			if err != nil {
				return api.FailedAuthentication
			}
			headers := c.Request().Header
			headers.Set("Goduit-Subject", identity.Subject)
			headers.Set("Goduit-Client-Username", identity.Username)
			headers.Set("Goduit-Client-Email", identity.UserEmail)
			return next(c)
		}
	}
}

func GenerateToken(userEmail, username, userID string) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, &Identity{
		UserEmail: userEmail,
		Username:  username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "goduit",
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(now.Add(10 * time.Minute)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        uuid.NewString(),
		},
	})
	return token.SignedString(encryptionkeys.PrivateKey)
}

func FromToken(authorizationHeader string) (*Identity, error) {
	authorizationHeader = strings.TrimPrefix(authorizationHeader, "Bearer ")
	token, err := jwt.ParseWithClaims(authorizationHeader, &Identity{}, func(t *jwt.Token) (interface{}, error) {
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
	claims, ok := token.Claims.(*Identity)
	if !ok {
		return nil, couldNotParseClaimsErr
	}
	return claims, nil
}
