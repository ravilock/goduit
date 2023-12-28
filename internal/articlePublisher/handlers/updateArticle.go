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

	request.Article.Slug = c.Param("slug")

	if err := request.Validate(); err != nil {
		return err
	}

	article := request.Model(authorUsername)

	ctx := c.Request().Context()

	currentArticle, err := h.service.GetArticleBySlug(ctx, request.Article.Slug)
	if err != nil {
		if appError := new(app.AppError); errors.As(err, &appError) {
			switch appError.ErrorCode {
			case app.ArticleNotFoundErrorCode:
				return api.ArticleNotFound(request.Article.Slug)
			}
		}
		return err
	}

	if authorUsername != *currentArticle.Author {
		return api.Forbidden
	}

	if err = h.service.UpdateArticle(ctx, request.Article.Slug, article); err != nil {
		if appError := new(app.AppError); errors.As(err, &appError) {
			switch appError.ErrorCode {
			case app.ArticleNotFoundErrorCode:
				return api.ArticleNotFound(request.Article.Slug)
			case app.ConflictErrorCode:
				return api.ConfictError
			}
		}
		return err
	}

	authorProfile, err := h.profileManager.GetProfileByUsername(ctx, authorUsername)
	if err != nil {
		if appError := new(app.AppError); errors.As(err, &appError) {
			switch appError.ErrorCode {
			case app.UserNotFoundErrorCode:
				return api.UserNotFound(authorUsername)
			}
		}
		return err
	}

	profileResponse, err := profileManagerAssembler.ProfileResponse(authorProfile, false)
	if err != nil {
		return err
	}

	response := assemblers.ArticleResponse(article, profileResponse)
	return c.JSON(http.StatusOK, response)
}
