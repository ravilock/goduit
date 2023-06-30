package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/api/assemblers"
	"github.com/ravilock/goduit/api/requests"
	"github.com/ravilock/goduit/api/validators"
	"github.com/ravilock/goduit/internal/app/services"
)

func Login(c echo.Context) error {
	request := new(requests.Login)
	if err := c.Bind(request); err != nil {
		return api.CouldNotUnmarshalBodyError
	}

	if err := validators.Login(request); err != nil {
		return err
	}

	dto := assemblers.Login(request)

	dto, err := services.Login(dto, c.Request().Context())
	if err != nil {
		return err
	}

	response := assemblers.Response(dto)

	return c.JSON(http.StatusOK, response)
}
