package handlers

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api/assemblers"
	"github.com/ravilock/goduit/internal/app/services"
)

func GetUser(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*jwt.RegisteredClaims)
	subject := claims.Subject

	dto := assemblers.GetUser(&subject)

	dto, err := services.GetUser(dto, c.Request().Context())
	if err != nil {
		return err
	}

	response := assemblers.Response(dto)

	return c.JSON(http.StatusOK, response)
}
