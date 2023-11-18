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

func UpdateUser(c echo.Context) error {
	request := new(requests.UpdateUser)
	if err := c.Bind(request); err != nil {
		return api.CouldNotUnmarshalBodyError
	}

	if err := validators.UpdateUser(request); err != nil {
		return err
	}

	model := request.Model()

	model, err := services.UpdateUser(model, request.User.Password, c.Request().Context())
	if err != nil {
		return err
	}

	response := assemblers.UserResponse(model, "")

	return c.JSON(http.StatusOK, response)
}
