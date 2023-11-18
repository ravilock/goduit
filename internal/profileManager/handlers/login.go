package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/internal/app"
	assemblers "github.com/ravilock/goduit/internal/profileManager/assembler"
	"github.com/ravilock/goduit/internal/profileManager/models"
	"github.com/ravilock/goduit/internal/profileManager/requests"
)

type logger interface {
	Login(ctx context.Context, model *models.User, password string) (*models.User, string, error)
}

type loginHandler struct {
	service logger
}

func (h *loginHandler) Login(c echo.Context) error {
	request := new(requests.Login)
	if err := c.Bind(request); err != nil {
		return api.CouldNotUnmarshalBodyError
	}

	if err := request.Validate(); err != nil {
		return err
	}

	user := request.Model()

	user, token, err := h.service.Login(c.Request().Context(), user, request.User.Password)
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
