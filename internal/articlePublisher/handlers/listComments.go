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
	profileManagerResponses "github.com/ravilock/goduit/internal/profileManager/responses"
)

type commentLister interface {
	ListComments(ctx context.Context, article string) ([]*models.Comment, error)
}

type listCommentsHandler struct {
	service          commentLister
	articlePublisher articleGetter
	profileManager   profileGetter
	followerCentral  isFollowedChecker
}

func (h *listCommentsHandler) ListComments(c echo.Context) error {
	request := new(requests.ArticleSlugRequest)
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

	ctx := c.Request().Context()

	article, err := h.articlePublisher.GetArticleBySlug(ctx, request.Slug)
	if err != nil {
		if appError := new(app.AppError); errors.As(err, &appError) {
			switch appError.ErrorCode {
			case app.ArticleNotFoundErrorCode:
				return api.ArticleNotFound(request.Slug)
			}
		}
		return err
	}

	comments, err := h.service.ListComments(ctx, article.ID.Hex())
	if err != nil {
		return err
	}

	authorMap := make(map[string]*profileManagerResponses.ProfileResponse)
	for _, comment := range comments {
		_, ok := authorMap[*comment.Author]
		if ok {
			continue
		}

		// TODO: Avaliar se eu devo apenas "pular" o comentário se o author não foi encontrado (ou atribuir para "deletado")
		author, err := h.profileManager.GetProfileByID(ctx, *comment.Author)
		if err != nil {
			if appError := new(app.AppError); errors.As(err, &appError) {
				switch appError.ErrorCode {
				case app.UserNotFoundErrorCode:
					return api.UserNotFound(identity.ClientUsername)
				}
			}
			return err
		}

		isFollowing := h.followerCentral.IsFollowedBy(ctx, author.ID.Hex(), identity.Subject)
		authorProfile, err := profileManagerAssembler.ProfileResponse(author, isFollowing)
		if err != nil {
			authorMap[*comment.Author] = nil
			continue
		}
		authorMap[*comment.Author] = authorProfile
	}

	response := new(responses.CommentsResponse)
	for _, comment := range comments {
		commentAuthor := authorMap[*comment.Author]
		commentResponse := assemblers.CommentResponse(comment, commentAuthor)
		response.Comment = append(response.Comment, commentResponse.Comment)
	}
	return c.JSON(http.StatusOK, response)
}
