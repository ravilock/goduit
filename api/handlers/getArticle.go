package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api/assemblers"
	"github.com/ravilock/goduit/api/requests"
	"github.com/ravilock/goduit/api/validators"
	"github.com/ravilock/goduit/internal/app/services"
)

func GetArticle(c echo.Context) error {
	username := c.Request().Header.Get("Goduit-Client-Username")
	request := new(requests.GetArticle)

	request.Slug = c.Param("slug")
	if err := validators.GetArticle(request); err != nil {
		return err
	}

	ctx := c.Request().Context()

	dto, err := services.GetArticleBySlug(request.Slug, username, ctx)
	if err != nil {
		return err
	}

	response := assemblers.ArticleResponse(dto)

	return c.JSON(http.StatusOK, response)
}