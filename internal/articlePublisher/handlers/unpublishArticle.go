package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/internal/app"
	"github.com/ravilock/goduit/internal/articlePublisher/requests"
	"github.com/ravilock/goduit/internal/identity"
)

type articleUnpublisher interface {
	UnpublishArticle(ctx context.Context, slug string) error
}

type UnpublishArticleHandler struct {
	articleUnpublisher articleUnpublisher
	articleGetter      articleGetter
}

func NewUnpublishArticleHandler(service articleUnpublisher, aarticleGetter articleGetter) *UnpublishArticleHandler {
	return &UnpublishArticleHandler{
		articleUnpublisher: service,
		articleGetter:      aarticleGetter,
	}
}

func (h *UnpublishArticleHandler) UnpublishArticle(c echo.Context) error {
	request := new(requests.ArticleSlugRequest)
	identity := new(identity.IdentityHeaders)
	binder := &echo.DefaultBinder{}
	if err := binder.BindPathParams(c, request); err != nil {
		return err
	}
	if err := binder.BindHeaders(c, identity); err != nil {
		return err
	}

	if err := request.Validate(); err != nil {
		return err
	}

	ctx := c.Request().Context()

	article, err := h.articleGetter.GetArticleBySlug(ctx, request.Slug)
	if err != nil {
		if appError := new(app.AppError); errors.As(err, &appError) {
			switch appError.ErrorCode {
			case app.ArticleNotFoundErrorCode:
				return api.ArticleNotFound(request.Slug)
			}
		}
		return err
	}

	if identity.Subject != *article.Author {
		return api.Forbidden
	}

	if err := h.articleUnpublisher.UnpublishArticle(ctx, request.Slug); err != nil {
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
