package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api/assemblers"
	"github.com/ravilock/goduit/api/requests"
	"github.com/ravilock/goduit/api/validators"
	"github.com/ravilock/goduit/internal/app/services"
)

func Unfollow(c echo.Context) error {
	clientUsername := c.Request().Header.Get("Goduit-Client-Username")
	request := new(requests.Follower)

	request.Username = c.Param("username")
	if err := validators.Follower(request); err != nil {
		return err
	}

	ctx := c.Request().Context()

	if err := services.Unfollow(request.Username, clientUsername, ctx); err != nil {
		return err
	}

	dto, err := services.GetProfileByUsername(request.Username, ctx)
	if err != nil {
		return err
	}
	dto.Following = false

	response, err := assemblers.ProfileResponse(dto)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response)
}
