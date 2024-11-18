package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/internal/app"
	"github.com/ravilock/goduit/internal/followerCentral/requests"
	"github.com/ravilock/goduit/internal/identity"
	"github.com/ravilock/goduit/internal/profileManager/assemblers"
)

type userUnfollower interface {
	Unfollow(ctx context.Context, followed, following string) error
}

type unfollowUserHandler struct {
	service        userUnfollower
	profileManager profileGetter
}

func (h *unfollowUserHandler) Unfollow(c echo.Context) error {
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

	fmt.Println(ctx, request.Username)
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

	err = h.service.Unfollow(ctx, followedUser.ID.Hex(), identity.Subject)
	if err != nil {
		return err
	}

	response, err := assemblers.ProfileResponse(followedUser, false)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response)
}
