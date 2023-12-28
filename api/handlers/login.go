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

func Login(c echo.Context) error {
	request := new(requests.Login)
	if err := c.Bind(request); err != nil {
		return api.CouldNotUnmarshalBodyError
	}

	if err := validators.Login(request); err != nil {
		return err
	}

	user := request.Model()

	user, token, err := services.Login(user, request.User.Password, c.Request().Context())
	if err != nil {
		if appError := new(app.AppError); errors.As(err, &appError) {
			switch appError.ErrorCode {
			case app.UserNotFoundErrorCode:
				fallthrough
			case app.WrongPasswordErrorCode:
				return api.FailedLoginAttempt
			}
		}
		return err
	}

	response := assemblers.UserResponse(user, token)

	return c.JSON(http.StatusOK, response)
}
