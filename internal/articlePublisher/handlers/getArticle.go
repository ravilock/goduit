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
	profileManagerModels "github.com/ravilock/goduit/internal/profileManager/models"
)

type articleGetter interface {
	GetArticleBySlug(ctx context.Context, slug string) (*models.Article, error)
}

type profileGetter interface {
	GetProfileByUsername(ctx context.Context, username string) (*profileManagerModels.User, error)
}

type isFollowedChecker interface {
	IsFollowedBy(ctx context.Context, followed, following string) bool
}

type getArticleHandler struct {
	service         articleGetter
	profileManager  profileGetter
	followerCentral isFollowedChecker
}

func (h *getArticleHandler) GetArticle(c echo.Context) error {
	clientUsername := c.Request().Header.Get("Goduit-Client-Username")
	request := new(requests.ArticleSlugRequest)
	if err := c.Bind(request); err != nil {
		return api.CouldNotUnmarshalBodyError
	}

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

	author, err := h.profileManager.GetProfileByUsername(ctx, *article.Author)
	if err != nil {
		return err
	}

	isFollowing := h.followerCentral.IsFollowedBy(ctx, *author.Username, clientUsername)

	authorProfile, err := profileManagerAssembler.ProfileResponse(author, isFollowing)
	if err != nil {
		return err
	}

	response := assemblers.ArticleResponse(article, authorProfile)

	return c.JSON(http.StatusOK, response)
}
