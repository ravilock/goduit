package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api/assemblers"
	"github.com/ravilock/goduit/internal/app/services"
)

func GetUser(c echo.Context) error {
	subject := c.Request().Header.Get("Goduit-Subject")

	model, err := services.GetUserByEmail(subject, c.Request().Context())
	if err != nil {
		return err
	}

	response := assemblers.UserResponse(model, nil)

	return c.JSON(http.StatusOK, response)
}
