package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/internal/app"
	"github.com/ravilock/goduit/internal/followerCentral/requests"
	"github.com/ravilock/goduit/internal/profileManager/assemblers"
	profileManager "github.com/ravilock/goduit/internal/profileManager/handlers"
)

type userFollower interface {
	Follow(ctx context.Context, followed, following string) error
}

type followUserHandler struct {
	service        userFollower
	profileManager profileManager.ProfileGetter
}

func (h *followUserHandler) Follow(c echo.Context) error {
	clientUsername := c.Request().Header.Get("Goduit-Client-Username")
	request := new(requests.Follower)

	request.Username = c.Param("username")
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

	err = h.service.Follow(ctx, request.Username, clientUsername)
	if err != nil {
		return err
	}

	response, err := assemblers.ProfileResponse(followedUser, true)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response)
}