package handlers

import (
	"context"
	"errors"
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
	logger
	UpdateProfile(ctx context.Context, subjectEmail, clientUsername, password string, model *models.User) error
}

type updateProfileHandler struct {
	service profileUpdater
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

	if err := h.service.UpdateProfile(c.Request().Context(), identityHeaders.ClientEmail, identityHeaders.ClientUsername, request.User.Password, model); err != nil {
		if appError := new(app.AppError); errors.As(err, &appError) {
			switch appError.ErrorCode {
			case app.ConflictErrorCode:
				return api.ConfictError
			}
		}
		return err
	}

	var tokenString string
	var err error
	if shouldGenerateNewToken(identityHeaders, model) {
		tokenString, err = identity.GenerateToken(*model.Email, *model.Username, model.ID.Hex())
		if err != nil {
			return err
		}
	}

	response := assemblers.UserResponse(model, tokenString)

	return c.JSON(http.StatusOK, response)
}

func shouldGenerateNewToken(identityHeaders *identity.IdentityHeaders, model *models.User) bool {
	return identityHeaders.ClientEmail != *model.Email || identityHeaders.ClientUsername != *model.Username
}
