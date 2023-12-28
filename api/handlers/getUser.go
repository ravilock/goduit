package handlers

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api/assemblers"
	"github.com/ravilock/goduit/internal/app/dtos"
	"github.com/ravilock/goduit/internal/app/services"
)

func GetUser(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*dtos.TokenClaims)
	subject := claims.Subject

	dto, err := services.GetUserByEmail(subject, c.Request().Context())
	if err != nil {
		return err
	}

	response := assemblers.UserResponse(dto)

	return c.JSON(http.StatusOK, response)
}
