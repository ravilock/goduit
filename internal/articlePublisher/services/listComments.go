package services

import (
	"context"

	"github.com/ravilock/goduit/internal/articlePublisher/models"
)

type commentLister interface {
	ListComments(ctx context.Context, article string) ([]*models.Comment, error)
}

type ListCommentsService struct {
	repository commentLister
}

func NewListCommentsService(repository commentLister) *ListCommentsService {
	return &ListCommentsService{
		repository: repository,
	}
}

func (s *ListCommentsService) ListComments(ctx context.Context, article string) ([]*models.Comment, error) {
	return s.repository.ListComments(ctx, article)
}
