package services

import (
	"context"

	"github.com/ravilock/goduit/internal/articlePublisher/models"
)

type commentGetter interface {
	GetCommentByID(ctx context.Context, ID string) (*models.Comment, error)
}

type GetCommentService struct {
	repository commentGetter
}

func NewGetCommentService(repository commentGetter) *GetCommentService {
	return &GetCommentService{
		repository: repository,
	}
}

func (s *GetCommentService) GetCommentByID(ctx context.Context, ID string) (*models.Comment, error) {
	return s.repository.GetCommentByID(ctx, ID)
}
