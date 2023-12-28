package handlers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/internal/profileManager/assemblers"
	"github.com/ravilock/goduit/internal/profileManager/models"
)

type ProfileGetter interface {
	GetProfileByEmail(ctx context.Context, email string) (*models.User, error)
	GetProfileByUsername(ctx context.Context, username string) (*models.User, error)
}

type getOwnProfileHandler struct {
	service ProfileGetter
}

func (h *getOwnProfileHandler) GetOwnProfile(c echo.Context) error {
	subject := c.Request().Header.Get("Goduit-Subject")

	user, err := h.service.GetProfileByEmail(c.Request().Context(), subject)
	if err != nil {
		return err
	}

	response := assemblers.UserResponse(user, "")

	return c.JSON(http.StatusOK, response)
}
