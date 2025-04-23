package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/internal/app"
	"github.com/ravilock/goduit/internal/articlePublisher/models"
	"github.com/ravilock/goduit/internal/articlePublisher/requests"
	"github.com/ravilock/goduit/internal/identity"
)

type commentDeleter interface {
	DeleteComment(ctx context.Context, ID string) error
	GetCommentByID(ctx context.Context, ID string) (*models.Comment, error)
}

type DeleteCommentHandler struct {
	service          commentDeleter
	articlePublisher articleGetter
}

func NewDeleteCommentHandler(
	service commentDeleter,
	articlePublisher articleGetter,
) *DeleteCommentHandler {
	return &DeleteCommentHandler{
		service:          service,
		articlePublisher: articlePublisher,
	}
}

func (h *DeleteCommentHandler) DeleteComment(c echo.Context) error {
	request := new(requests.DeleteCommentRequest)
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

	_, err := h.articlePublisher.GetArticleBySlug(ctx, request.Slug)
	if err != nil {
		if appError := new(app.AppError); errors.As(err, &appError) {
			switch appError.ErrorCode {
			case app.ArticleNotFoundErrorCode:
				return api.ArticleNotFound(request.Slug)
			}
		}
		return err
	}

	comment, err := h.service.GetCommentByID(ctx, request.ID)
	if err != nil {
		if appError := new(app.AppError); errors.As(err, &appError) {
			switch appError.ErrorCode {
			case app.CommentNotFoundErrorCode:
				return api.CommentNotFound(request.ID)
			}
		}
		return err
	}

	if identity.Subject != *comment.Author {
		return api.Forbidden
	}

	if err := h.service.DeleteComment(ctx, request.ID); err != nil {
		if appError := new(app.AppError); errors.As(err, &appError) {
			switch appError.ErrorCode {
			case app.CommentNotFoundErrorCode:
				return api.CommentNotFound(request.ID)
			}
		}
		return err
	}
	return c.NoContent(http.StatusNoContent)
}
