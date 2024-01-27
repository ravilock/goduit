package services

import (
	"context"

	"github.com/ravilock/goduit/internal/articlePublisher/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type commentLister interface {
	ListComments(ctx context.Context, article primitive.ObjectID) ([]*models.Comment, error)
}

type listCommentService struct {
	repository commentLister
}

func (s *listCommentService) ListComments(ctx context.Context, article primitive.ObjectID) ([]*models.Comment, error) {
	return s.repository.ListComments(ctx, article)
}
