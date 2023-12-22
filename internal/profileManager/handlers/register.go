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

type profileRegister interface {
	Register(ctx context.Context, model *models.User, password string) (*models.User, string, error)
}

type registerProfileHandler struct {
	service profileRegister
}

func (h *registerProfileHandler) Register(c echo.Context) error {
	request := new(requests.Register)
	if err := c.Bind(request); err != nil {
		return api.CouldNotUnmarshalBodyError
	}

	if err := request.Validate(); err != nil {
		return err
	}

	user := request.Model()

	user, token, err := h.service.Register(c.Request().Context(), user, request.User.Password)
	if err != nil {
		if appError := new(app.AppError); errors.As(err, &appError) {
			switch appError.ErrorCode {
			case app.ConflictErrorCode:
				return api.ConfictError
			}
		}
		return err
	}

	response := assemblers.UserResponse(user, token)

	return c.JSON(http.StatusCreated, response)
}
