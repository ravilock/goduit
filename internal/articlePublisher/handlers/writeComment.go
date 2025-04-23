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

type commentWriter interface {
	WriteComment(ctx context.Context, comment *models.Comment) error
}

type WriteCommentHandler struct {
	service          commentWriter
	articlePublisher articleGetter
	profileManager   profileGetter
}

func NewWriteCommentHandler(
	service commentWriter,
	articlePublisher articleGetter,
	profileManager profileGetter,
) *WriteCommentHandler {
	return &WriteCommentHandler{
		service:          service,
		articlePublisher: articlePublisher,
		profileManager:   profileManager,
	}
}

func (h *WriteCommentHandler) WriteComment(c echo.Context) error {
	request := new(requests.WriteCommentRequest)
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

	comment := request.Model(identity.Subject)

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
	articleID := article.ID.Hex()
	comment.Article = &articleID

	if err := h.service.WriteComment(ctx, comment); err != nil {
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

	response := assemblers.CommentResponse(comment, profileResponse)
	return c.JSON(http.StatusCreated, response)
}
