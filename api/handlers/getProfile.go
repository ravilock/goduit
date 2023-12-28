package handlers

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api/assemblers"
	"github.com/ravilock/goduit/api/requests"
	"github.com/ravilock/goduit/api/validators"
	"github.com/ravilock/goduit/internal/app/dtos"
	"github.com/ravilock/goduit/internal/app/services"
)

func GetProfile(c echo.Context) error {
	var subject string
	request := new(requests.GetProfile)

	request.Username = c.Param("username")
	if err := validators.GetProfile(request); err != nil {
		return err
	}

	tokenString := c.Get("user")
	if tokenString != "" {
		claims := tokenString.(*jwt.Token).Claims.(*dtos.TokenClaims)
		if request.Username != claims.Username {
			subject = claims.Username
		}
	}

	dto, err := services.GetProfileByUsername(request.Username, subject, c.Request().Context())
	if err != nil {
		return err
	}

	response := assemblers.ProfileResponse(dto)

	return c.JSON(http.StatusOK, response)
}
