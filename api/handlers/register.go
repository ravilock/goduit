package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/api/assemblers"
	"github.com/ravilock/goduit/api/requests"
	"github.com/ravilock/goduit/api/validators"
	"github.com/ravilock/goduit/internal/app/services"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register(c echo.Context) error {
	request := new(requests.Register)
	if err := c.Bind(request); err != nil {
		return api.CouldNotUnmarshalBodyError
	}

	if err := validators.Register(request); err != nil {
		return err
	}

	dto := assemblers.Register(request)

	dto, err := services.Register(dto, c.Request().Context())
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return api.ConfictError
		}
		return err
	}

	response := assemblers.UserResponse(dto)

	return c.JSON(http.StatusCreated, response)
}
