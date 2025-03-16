package handlers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/internal/identity"
	"github.com/ravilock/goduit/internal/profileManager/assemblers"
	"github.com/ravilock/goduit/internal/profileManager/models"
)

type profileGetter interface {
	GetProfileByUsername(ctx context.Context, username string) (*models.User, error)
	GetProfileByID(ctx context.Context, ID string) (*models.User, error)
}

type getOwnProfileHandler struct {
	service profileGetter
}

func (h *getOwnProfileHandler) GetOwnProfile(c echo.Context) error {
	identity := new(identity.IdentityHeaders)
	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, identity); err != nil {
		return err
	}

	user, err := h.service.GetProfileByID(c.Request().Context(), identity.Subject)
	if err != nil {
		return err
	}

	response := assemblers.UserResponse(user)

	return c.JSON(http.StatusOK, response)
}
