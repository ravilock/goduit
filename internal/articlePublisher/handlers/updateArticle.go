package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/internal/app"
	"github.com/ravilock/goduit/internal/articlePublisher/assemblers"
	"github.com/ravilock/goduit/internal/articlePublisher/models"
	"github.com/ravilock/goduit/internal/articlePublisher/requests"
	profileManagerAssembler "github.com/ravilock/goduit/internal/profileManager/assemblers"
)

type articleUpdater interface {
	UpdateArticle(ctx context.Context, slug string, article *models.Article) error
	articleGetter
}

type updateArticleHandler struct {
	service        articleUpdater
	profileManager profileGetter
}

func (h *updateArticleHandler) UpdateArticle(c echo.Context) error {
	authorUsername := c.Request().Header.Get("Goduit-Client-Username")
	request := new(requests.UpdateArticle)
	if err := c.Bind(request); err != nil {
		return api.CouldNotUnmarshalBodyError
	}

	currentSlug := c.Param("slug")

	if err := request.Validate(); err != nil {
		return err
	}

	article := request.Model(authorUsername)

	ctx := c.Request().Context()

	article, err := h.service.GetArticleBySlug(ctx, currentSlug)
	if err != nil {
		if appError := new(app.AppError); errors.As(err, &appError) {
			switch appError.ErrorCode {
			case app.ArticleNotFoundErrorCode:
				return api.ArticleNotFound(currentSlug)
			}
		}
		return err
	}

	authorProfile, err := h.profileManager.GetProfileByUsername(ctx, authorUsername)
	if err != nil {
		return api.UserNotFound(authorUsername)
	}

	profileResponse, err := profileManagerAssembler.ProfileResponse(authorProfile, false)
	if err != nil {
		return err
	}

	response := assemblers.ArticleResponse(article, profileResponse)
	return c.JSON(http.StatusCreated, response)
}
