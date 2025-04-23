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
	"github.com/ravilock/goduit/internal/articlePublisher/responses"
	"github.com/ravilock/goduit/internal/identity"
	profileManagerAssembler "github.com/ravilock/goduit/internal/profileManager/assemblers"
)

type articleFeeder interface {
	FeedArticles(ctx context.Context, user string, limit, offset int64) ([]*models.Article, error)
}

type FeedArticlesHandler struct {
	service        articleFeeder
	profileManager profileGetter
}

func NewFeedArticlesHandler(service articleFeeder, profileManager profileGetter) *FeedArticlesHandler {
	return &FeedArticlesHandler{
		service:        service,
		profileManager: profileManager,
	}
}

func (h *FeedArticlesHandler) FeedArticles(c echo.Context) error {
	request := requests.NewFeedArticlesRequest()
	identity := new(identity.IdentityHeaders)
	binder := &echo.DefaultBinder{}
	if err := binder.BindQueryParams(c, request); err != nil {
		return err
	}
	if err := binder.BindHeaders(c, identity); err != nil {
		return err
	}

	if err := request.Validate(); err != nil {
		return err
	}

	ctx := c.Request().Context()

	articles, err := h.service.FeedArticles(ctx, identity.Subject, int64(request.Pagination.Limit), int64(request.Pagination.Offset))
	if err != nil {
		if appError := new(app.AppError); errors.As(err, &appError) {
			switch appError.ErrorCode {
			case app.FeedNotFoundErrorCode:
				return api.FeedNotFound(identity.Subject)
			}
		}
		return err
	}

	response := responses.ArticlesResponse{Articles: make([]responses.MultiArticle, 0, len(articles))}
	for _, article := range articles {
		// TODO: refactor so that if multiple articles from the same author and are in the same page, this loop wont repeat for each article
		author, err := h.profileManager.GetProfileByID(ctx, *article.Author)
		if err != nil {
			continue
		}

		authorProfile, err := profileManagerAssembler.ProfileResponse(author, true)
		if err != nil {
			continue
		}

		response.Articles = append(response.Articles, *assemblers.MultiArticleResponse(article, authorProfile))
	}

	return c.JSON(http.StatusOK, response)
}
