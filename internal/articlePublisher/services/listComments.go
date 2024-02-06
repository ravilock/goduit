package services

import (
	"context"

	"github.com/ravilock/goduit/internal/articlePublisher/models"
)

type commentLister interface {
	ListComments(ctx context.Context, article string) ([]*models.Comment, error)
}

type listCommentService struct {
	repository commentLister
}

func (s *listCommentService) ListComments(ctx context.Context, article string) ([]*models.Comment, error) {
	return s.repository.ListComments(ctx, article)
}
