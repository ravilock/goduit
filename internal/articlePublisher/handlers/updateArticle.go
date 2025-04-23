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
	"github.com/ravilock/goduit/internal/identity"
	profileManagerAssembler "github.com/ravilock/goduit/internal/profileManager/assemblers"
)

type articleUpdater interface {
	UpdateArticle(ctx context.Context, slug string, article *models.Article) error
}

type UpdateArticleHandler struct {
	articleUpdater articleUpdater
	articleGetter  articleGetter
	profileManager profileGetter
}

func NewUpdateArticleHandler(service articleUpdater, articleGetter articleGetter, profileManager profileGetter) *UpdateArticleHandler {
	return &UpdateArticleHandler{
		articleUpdater: service,
		articleGetter:  articleGetter,
		profileManager: profileManager,
	}
}

func (h *UpdateArticleHandler) UpdateArticle(c echo.Context) error {
	request := new(requests.UpdateArticleRequest)
	identity := new(identity.IdentityHeaders)
	binder := &echo.DefaultBinder{}
	if err := binder.BindBody(c, request); err != nil {
		return api.CouldNotUnmarshalBodyError
	}
	if err := binder.BindPathParams(c, request); err != nil {
		return err
	}
	if err := binder.BindHeaders(c, identity); err != nil {
		return err
	}

	if err := request.Validate(); err != nil {
		return err
	}

	article := request.Model()

	ctx := c.Request().Context()

	currentArticle, err := h.articleGetter.GetArticleBySlug(ctx, request.Slug)
	if err != nil {
		if appError := new(app.AppError); errors.As(err, &appError) {
			switch appError.ErrorCode {
			case app.ArticleNotFoundErrorCode:
				return api.ArticleNotFound(request.Slug)
			}
		}
		return err
	}

	if identity.Subject != *currentArticle.Author {
		return api.Forbidden
	}

	if err = h.articleUpdater.UpdateArticle(ctx, request.Slug, article); err != nil {
		if appError := new(app.AppError); errors.As(err, &appError) {
			switch appError.ErrorCode {
			case app.ArticleNotFoundErrorCode:
				return api.ArticleNotFound(request.Slug)
			case app.ConflictErrorCode:
				return api.ConfictError
			}
		}
		return err
	}

	authorProfile, err := h.profileManager.GetProfileByID(ctx, identity.Subject)
	if err != nil {
		if appError := new(app.AppError); errors.As(err, &appError) {
			switch appError.ErrorCode {
			case app.UserNotFoundErrorCode:
				return api.UserNotFound(identity.ClientUsername)
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
