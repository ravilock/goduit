package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/internal/app"
	"github.com/ravilock/goduit/internal/profileManager/assemblers"
	"github.com/ravilock/goduit/internal/profileManager/requests"
)

type IsFollowedChecker interface {
	IsFollowedBy(ctx context.Context, followed, following string) bool
}

type getProfileHandler struct {
	service         ProfileGetter
	followerCentral IsFollowedChecker
}

func (h *getProfileHandler) GetProfile(c echo.Context) error {
	clientUsername := c.Request().Header.Get("Goduit-Client-Username")
	fmt.Println(clientUsername)
	request := new(requests.GetProfile)

	request.Username = c.Param("username")
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

	isFollowing := h.followerCentral.IsFollowedBy(ctx, request.Username, clientUsername)

	response, err := assemblers.ProfileResponse(profile, isFollowing)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response)
}
