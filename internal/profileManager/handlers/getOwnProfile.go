package handlers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	assemblers "github.com/ravilock/goduit/internal/profileManager/assembler"
	"github.com/ravilock/goduit/internal/profileManager/models"
)

type profileGetter interface {
	GetProfile(ctx context.Context, email string) (*models.User, error)
}

type getOwnProfileHandler struct {
	service profileGetter
}

func (h *getOwnProfileHandler) GetOwnProfile(c echo.Context) error {
	subject := c.Request().Header.Get("Goduit-Subject")

	user, err := h.service.GetProfile(c.Request().Context(), subject)
	if err != nil {
		return err
	}

	response := assemblers.UserResponse(user, "")

	return c.JSON(http.StatusOK, response)
}
