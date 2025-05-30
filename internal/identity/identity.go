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
	"github.com/ravilock/goduit/internal/config"
	"github.com/ravilock/goduit/internal/cookie"
)

var (
	errInvalidToken       = errors.New("invalid Token")
	errCouldNotParseClaim = errors.New("could Not Parse Claims")
)

type Identity struct {
	UserEmail string `json:"userId,omitempty"`
	Username  string `json:"username,omitempty"`
	jwt.RegisteredClaims
}

type IdentityHeaders struct {
	Subject        string `header:"Goduit-Subject"`
	ClientUsername string `header:"Goduit-Client-Username"`
	ClientEmail    string `header:"Goduit-Client-Email"`
}

func CreateAuthMiddleware(requiredAuthentication bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			now := time.Now()
			cookie, err := c.Cookie(cookie.CookieKey)
			if err != nil {
				if requiredAuthentication {
					return api.FailedAuthentication
				}
				return next(c)
			}
			if cookie.Expires.After(now) {
				return api.FailedAuthentication
			}
			token := cookie.Value
			if !requiredAuthentication && token == "" {
				return next(c)
			}
			identity, err := FromToken(token)
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
	now := time.Now().UTC().Truncate(time.Millisecond)
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, &Identity{
		UserEmail: userEmail,
		Username:  username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "goduit",
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        uuid.NewString(),
		},
	})
	return token.SignedString(config.PrivateKey)
}

func FromToken(tokenString string) (*Identity, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	token, err := jwt.ParseWithClaims(tokenString, &Identity{}, func(t *jwt.Token) (interface{}, error) {
		key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(os.Getenv("PUBLIC_KEY")))
		if err != nil {
			return nil, err
		}
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errInvalidToken
		}
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errInvalidToken
	}
	claims, ok := token.Claims.(*Identity)
	if !ok {
		return nil, errCouldNotParseClaim
	}
	return claims, nil
}
