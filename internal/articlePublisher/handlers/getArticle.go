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
	profileManagerModels "github.com/ravilock/goduit/internal/profileManager/models"
)

type articleGetter interface {
	GetArticleBySlug(ctx context.Context, slug string) (*models.Article, error)
}

type profileGetter interface {
	GetProfileByID(ctx context.Context, ID string) (*profileManagerModels.User, error)
	GetProfileByUsername(ctx context.Context, username string) (*profileManagerModels.User, error)
}

type isFollowedChecker interface {
	IsFollowedBy(ctx context.Context, followed, following string) bool
}

type GetArticleHandler struct {
	service         articleGetter
	profileManager  profileGetter
	followerCentral isFollowedChecker
}

func NewGetArticleHandler(
	service articleGetter,
	profileManager profileGetter,
	followerCentral isFollowedChecker,
) *GetArticleHandler {
	return &GetArticleHandler{
		service:         service,
		profileManager:  profileManager,
		followerCentral: followerCentral,
	}
}

func (h *GetArticleHandler) GetArticle(c echo.Context) error {
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

	author, err := h.profileManager.GetProfileByID(ctx, *article.Author)
	if err != nil {
		return err
	}

	isFollowing := h.followerCentral.IsFollowedBy(ctx, author.ID.Hex(), identity.Subject)

	authorProfile, err := profileManagerAssembler.ProfileResponse(author, isFollowing)
	if err != nil {
		return err
	}

	response := assemblers.ArticleResponse(article, authorProfile)

	return c.JSON(http.StatusOK, response)
}
