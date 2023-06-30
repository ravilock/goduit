package middlewares

import (
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
)

func AuthenticationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(os.Getenv("PUBLIC_KEY")))
		if err != nil {
			return err
		}

		token, err := jwt.ParseWithClaims(c.Request().Header.Get("Authorization"), &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
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

		_, ok := token.Claims.(*jwt.RegisteredClaims)
		if !ok {
			log.Println("Could Not Parse Claims")
			return api.FailedAuthentication
		}
		return nil
	}
}
