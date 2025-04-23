package services

import (
	"context"

	"github.com/ravilock/goduit/internal/articlePublisher/models"
)

type commentWriter interface {
	WriteComment(ctx context.Context, article *models.Comment) error
}

type WriteCommentService struct {
	repository commentWriter
}

func NewWriteCommentService(repository commentWriter) *WriteCommentService {
	return &WriteCommentService{
		repository: repository,
	}
}

func (s *WriteCommentService) WriteComment(ctx context.Context, comment *models.Comment) error {
	if err := s.repository.WriteComment(ctx, comment); err != nil {
		return err
	}
	return nil
}
