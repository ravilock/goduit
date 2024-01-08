package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/internal/app"
	"github.com/ravilock/goduit/internal/followerCentral/requests"
	"github.com/ravilock/goduit/internal/identity"
	"github.com/ravilock/goduit/internal/profileManager/assemblers"
	profileManagerModels "github.com/ravilock/goduit/internal/profileManager/models"
)

type userFollower interface {
	Follow(ctx context.Context, followed, following string) error
}

type profileGetter interface {
	GetProfileByUsername(ctx context.Context, username string) (*profileManagerModels.User, error)
}

type followUserHandler struct {
	service        userFollower
	profileManager profileGetter
}

func (h *followUserHandler) Follow(c echo.Context) error {
	request := new(requests.FollowerRequest)
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

	followedUser, err := h.profileManager.GetProfileByUsername(ctx, request.Username)
	if err != nil {
		if appError := new(app.AppError); errors.As(err, &appError) {
			switch appError.ErrorCode {
			case app.UserNotFoundErrorCode:
				return api.UserNotFound(request.Username)
			}
		}
		return err
	}

	err = h.service.Follow(ctx, followedUser.ID.Hex(), identity.Subject)
	if err != nil {
		return err
	}

	response, err := assemblers.ProfileResponse(followedUser, true)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response)
}
