package services

import (
	"context"

	"github.com/ravilock/goduit/internal/articlePublisher/models"
)

type commentGetter interface {
	GetCommentByID(ctx context.Context, ID string) (*models.Comment, error)
}

type getCommentService struct {
	repository commentGetter
}

func (s *getCommentService) GetCommentByID(ctx context.Context, ID string) (*models.Comment, error) {
	return s.repository.GetCommentByID(ctx, ID)
}
