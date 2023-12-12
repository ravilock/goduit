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
	UpdateProfile(ctx context.Context, subjectEmail, clientUsername, password string, model *models.User) (*models.User, string, error)
}

type updateProfileHandler struct {
	service profileUpdater
}

func (h *updateProfileHandler) UpdateProfile(c echo.Context) error {
	subjectEmail := c.Request().Header.Get("Goduit-Subject")
	clientUsername := c.Request().Header.Get("Goduit-Client-Username")

	request := new(requests.UpdateProfile)
	if err := c.Bind(request); err != nil {
		return api.CouldNotUnmarshalBodyError
	}

	if err := request.Validate(); err != nil {
		return err
	}

	model := request.Model()

	model, token, err := h.service.UpdateProfile(c.Request().Context(), subjectEmail, clientUsername, request.User.Password, model)
	if err != nil {
		return err
	}

	response := assemblers.UserResponse(model, token)

	return c.JSON(http.StatusOK, response)
}
