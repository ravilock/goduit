package middlewares

import (
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/internal/app/dtos"
)

func CreateAuthMiddleware(requiredAuthentication bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authorizationHeader := c.Request().Header.Get("Authorization")
			if !requiredAuthentication && authorizationHeader == "" {
				c.Set("user", "")
				return next(c)
			}

			key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(os.Getenv("PUBLIC_KEY")))
			if err != nil {
				return err
			}

			token, err := jwt.ParseWithClaims(authorizationHeader, &dtos.TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, api.UnexpectedTokenSigningMethod(t.Method.Alg())
				}
				return key, nil
			})
			if err != nil {
				log.Println("Could Not Parse Token, err=", err)
				return api.FailedAuthentication
			}
			if !token.Valid {
				log.Println("Invalid Token")
				return api.FailedAuthentication
			}

			_, ok := token.Claims.(*dtos.TokenClaims)
			if !ok {
				log.Println("Could Not Parse Claims")
				return api.FailedAuthentication
			}
			c.Set("user", token)
			return next(c)
		}
	}
}
