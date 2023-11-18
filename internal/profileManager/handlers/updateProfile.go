package handlers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/internal/profileManager/assemblers"
	"github.com/ravilock/goduit/internal/profileManager/models"
	"github.com/ravilock/goduit/internal/profileManager/requests"
)

type profileUpdater interface {
	UpdateProfile(ctx context.Context, model *models.User, password string) (*models.User, error)
}

type updateProfileHandler struct {
	service profileUpdater
}

func (h *updateProfileHandler) UpdateProfile(c echo.Context) error {
	request := new(requests.UpdateProfile)
	if err := c.Bind(request); err != nil {
		return api.CouldNotUnmarshalBodyError
	}

	if err := request.Validate(); err != nil {
		return err
	}

	model := request.Model()

	model, err := h.service.UpdateProfile(c.Request().Context(), model, request.User.Password)
	if err != nil {
		return err
	}

	response := assemblers.UserResponse(model, "")

	return c.JSON(http.StatusOK, response)
}
