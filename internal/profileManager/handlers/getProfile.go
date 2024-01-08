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
	"github.com/ravilock/goduit/internal/profileManager/requests"
)

type isFollowedChecker interface {
	IsFollowedBy(ctx context.Context, followed, following string) bool
}

type getProfileHandler struct {
	service         profileGetter
	followerCentral isFollowedChecker
}

func (h *getProfileHandler) GetProfile(c echo.Context) error {
	request := new(requests.GetProfileRequest)
	identity := new(identity.IdentityHeaders)
	binder := &echo.DefaultBinder{}
	if err := binder.BindPathParams(c, request); err != nil {
		return api.CouldNotUnmarshalBodyError
	}
	if err := binder.BindHeaders(c, identity); err != nil {
		return err
	}

	if err := request.Validate(); err != nil {
		return err
	}

	ctx := c.Request().Context()
	profile, err := h.service.GetProfileByUsername(ctx, request.Username)
	if err != nil {
		if appError := new(app.AppError); errors.As(err, &appError) {
			switch appError.ErrorCode {
			case app.UserNotFoundErrorCode:
				return api.UserNotFound(request.Username)
			}
		}
		return err
	}

	isFollowing := h.followerCentral.IsFollowedBy(ctx, profile.ID.Hex(), identity.Subject)

	response, err := assemblers.ProfileResponse(profile, isFollowing)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response)
}
