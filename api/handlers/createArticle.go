package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/api/assemblers"
	"github.com/ravilock/goduit/api/requests"
	"github.com/ravilock/goduit/api/validators"
	"github.com/ravilock/goduit/internal/app/services"
)

func CreateArticle(c echo.Context) error {
	username := c.Request().Header.Get("Goduit-Client-Username")

	request := new(requests.CreateArticle)
	if err := c.Bind(request); err != nil {
		return api.CouldNotUnmarshalBodyError
	}

	if err := validators.CreateArticle(request); err != nil {
		return err
	}

	ctx := c.Request().Context()

	author, err := services.GetProfileByUsername(username, ctx)
	if err != nil {
		return err
	}

	authorProfile, err := assemblers.ProfileResponse(author, false)
	if err != nil {
		return err
	}

	model := request.Model(author.Username)

	model, err = services.CreateArticle(model, ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, assemblers.ArticleResponse(model, *authorProfile))
}
