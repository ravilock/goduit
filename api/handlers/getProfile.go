package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api/assemblers"
	"github.com/ravilock/goduit/api/requests"
	"github.com/ravilock/goduit/api/validators"
	"github.com/ravilock/goduit/internal/app/services"
)

func GetProfile(c echo.Context) error {
	clientUsername := c.Request().Header.Get("Goduit-Client-Username")
	request := new(requests.GetProfile)

	request.Username = c.Param("username")
	if err := validators.GetProfile(request); err != nil {
		return err
	}

	ctx := c.Request().Context()

	dto, err := services.GetProfileByUsername(request.Username, ctx)
	if err != nil {
		return err
	}

	dto.Following = services.IsFollowedBy(request.Username, clientUsername, ctx)

	response, err := assemblers.ProfileResponse(dto)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response)
}