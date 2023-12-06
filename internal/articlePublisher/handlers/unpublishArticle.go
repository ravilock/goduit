package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/internal/app"
	"github.com/ravilock/goduit/internal/articlePublisher/requests"
)

type articleUnpublisher interface {
	UnpublishArticle(ctx context.Context, slug string) error
	articleGetter
}

type unpublishArticleHandler struct {
	service articleUnpublisher
}

func (h *unpublishArticleHandler) UnpublishArticle(c echo.Context) error {
	clientUsername := c.Request().Header.Get("Goduit-Client-Username")
	request := new(requests.UnpublishArticle)

	request.Slug = c.Param("slug")
	if err := request.Validate(); err != nil {
		return err
	}

	ctx := c.Request().Context()

	article, err := h.service.GetArticleBySlug(ctx, request.Slug)
	if err != nil {
		if appError := new(app.AppError); errors.As(err, &appError) {
			switch appError.ErrorCode {
			case app.ArticleNotFoundErrorCode:
				return api.ArticleNotFound(request.Slug)
			}
		}
		return err
	}

	if *article.Author != clientUsername {
		return api.Forbidden
	}

	if err := h.service.UnpublishArticle(ctx, request.Slug); err != nil {
		if appError := new(app.AppError); errors.As(err, &appError) {
			switch appError.ErrorCode {
			case app.ArticleNotFoundErrorCode:
				return api.ArticleNotFound(request.Slug)
			}
		}
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
