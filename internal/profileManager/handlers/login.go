package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/internal/app"
	"github.com/ravilock/goduit/internal/profileManager/assemblers"
	"github.com/ravilock/goduit/internal/profileManager/models"
	"github.com/ravilock/goduit/internal/profileManager/requests"
)

type logger interface {
	Login(ctx context.Context, email, password string) (*models.User, string, error)
}

type loginHandler struct {
	service logger
}

func (h *loginHandler) Login(c echo.Context) error {
	request := new(requests.LoginRequest)
	if err := c.Bind(request); err != nil {
		return api.CouldNotUnmarshalBodyError
	}

	if err := request.Validate(); err != nil {
		return err
	}

	user, token, err := h.service.Login(c.Request().Context(), request.User.Email, request.User.Password)
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
