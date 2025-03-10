package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/internal/app"
	"github.com/ravilock/goduit/internal/articlePublisher/assemblers"
	"github.com/ravilock/goduit/internal/articlePublisher/models"
	"github.com/ravilock/goduit/internal/articlePublisher/requests"
	"github.com/ravilock/goduit/internal/articlePublisher/responses"
	"github.com/ravilock/goduit/internal/identity"
	profileManagerAssembler "github.com/ravilock/goduit/internal/profileManager/assemblers"
)

type articleLister interface {
	ListArticles(ctx context.Context, author, tag string, limit, offset int64) ([]*models.Article, error)
}

type listArticlesHandler struct {
	service         articleLister
	profileManager  profileGetter
	followerCentral isFollowedChecker
}

func (h *listArticlesHandler) ListArticles(c echo.Context) error {
	request := requests.NewListArticlesRequest()
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

	if request.Filters.Author != "" {
		author, err := h.profileManager.GetProfileByUsername(ctx, request.Filters.Author)
		if err != nil {
			if appError := new(app.AppError); !errors.As(err, &appError) {
				return err
			}
		} else {
			request.Filters.Author = author.ID.Hex()
		}
	}

	articles, err := h.service.ListArticles(ctx, request.Filters.Author, request.Filters.Tag, int64(request.Pagination.Limit), int64(request.Pagination.Offset))
	if err != nil {
		return err
	}

	response := responses.ArticlesResponse{Articles: make([]responses.MultiArticle, 0, len(articles))}
	for _, article := range articles {
		// TODO: refactor so that if multiple articles from the same author and are in the same page, this loop wont repeat for each article
		author, err := h.profileManager.GetProfileByID(ctx, *article.Author)
		if err != nil {
			continue
		}

		isFollowing := h.followerCentral.IsFollowedBy(ctx, *article.Author, identity.Subject)

		authorProfile, err := profileManagerAssembler.ProfileResponse(author, isFollowing)
		if err != nil {
			continue
		}

		response.Articles = append(response.Articles, *assemblers.MultiArticleResponse(article, authorProfile))
	}

	return c.JSON(http.StatusOK, response)
}
