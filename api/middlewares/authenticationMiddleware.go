package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/internal/identity"
)

func CreateAuthMiddleware(requiredAuthentication bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authorizationHeader := c.Request().Header.Get("Authorization")
			if !requiredAuthentication && authorizationHeader == "" {
				c.Set("user", "")
				return next(c)
			}

			identity, err := identity.FromToken(authorizationHeader)
			if err != nil {
				return api.FailedAuthentication
			}
			headers := c.Request().Header
			headers.Set("Goduit-Client-Username", identity.Username)
			headers.Set("Goduit-Subject", identity.Subject)
			return next(c)
		}
	}
}
