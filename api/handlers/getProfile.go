package handlers

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api/assemblers"
	"github.com/ravilock/goduit/internal/app/dtos"
	"github.com/ravilock/goduit/internal/app/services"
)

func GetProfile(c echo.Context) error {
	var subject string
	username, tokenString := c.Param("username"), c.Get("user")
	if tokenString != "" {
		claims := tokenString.(*jwt.Token).Claims.(*dtos.TokenClaims)
		if username != claims.Username {
			subject = claims.Username
		}
	}

	dto, err := services.GetProfileByUsername(username, subject, c.Request().Context())
	if err != nil {
		return err
	}

	response := assemblers.ProfileResponse(dto)

	return c.JSON(http.StatusOK, response)
}
