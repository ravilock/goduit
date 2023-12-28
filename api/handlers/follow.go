package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/api/assemblers"
	"github.com/ravilock/goduit/api/requests"
	"github.com/ravilock/goduit/api/validators"
	"github.com/ravilock/goduit/internal/app"
	"github.com/ravilock/goduit/internal/app/services"
)

func Follow(c echo.Context) error {
	clientUsername := c.Request().Header.Get("Goduit-Client-Username")
	request := new(requests.Follower)

	request.Username = c.Param("username")
	if err := validators.Follower(request); err != nil {
		return err
	}

	ctx := c.Request().Context()

	err := services.Follow(request.Username, clientUsername, ctx)

	model, err := services.GetProfileByUsername(request.Username, ctx)
	if err != nil {
		if appError := new(app.AppError); errors.As(err, &appError) {
			switch appError.ErrorCode {
			case app.UserNotFoundErrorCode:
				return api.UserNotFound(request.Username)
			}
		}
		return err
	}

	response, err := assemblers.ProfileResponse(model, true)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response)
}
