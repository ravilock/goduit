package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/internal/app"
	"github.com/ravilock/goduit/internal/identity"
	"github.com/ravilock/goduit/internal/profileManager/assemblers"
	"github.com/ravilock/goduit/internal/profileManager/models"
	"github.com/ravilock/goduit/internal/profileManager/requests"
)

type profileUpdater interface {
	UpdateProfile(ctx context.Context, subjectEmail, clientUsername, password string, model *models.User) (string, error)
}

type updateProfileHandler struct {
	service       profileUpdater
	cookieService CookieCreator
}

func (h *updateProfileHandler) UpdateProfile(c echo.Context) error {
	request := new(requests.UpdateProfileRequest)
	identityHeaders := new(identity.IdentityHeaders)
	binder := &echo.DefaultBinder{}
	if err := binder.BindBody(c, request); err != nil {
		return api.CouldNotUnmarshalBodyError
	}
	if err := binder.BindHeaders(c, identityHeaders); err != nil {
		return err
	}

	if err := request.Validate(); err != nil {
		return err
	}

	model := request.Model()

	token, err := h.service.UpdateProfile(c.Request().Context(), identityHeaders.ClientEmail, identityHeaders.ClientUsername, request.User.Password, model)
	if err != nil {
		if appError := new(app.AppError); errors.As(err, &appError) {
			switch appError.ErrorCode {
			case app.ConflictErrorCode:
				return api.ConfictError
			case app.UserNotFoundErrorCode:
				return api.UserNotFound(fmt.Sprintf("%s+%s", identityHeaders.ClientEmail, identityHeaders.ClientUsername))
			}
		}
		return err
	}

	response := assemblers.UserResponse(model)
	if token != "" {
		cookie := h.cookieService.Create(token)
		c.SetCookie(cookie)
	}
	return c.JSON(http.StatusOK, response)
}
