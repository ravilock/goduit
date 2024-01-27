package services

import (
	"context"

	"github.com/ravilock/goduit/internal/articlePublisher/models"
)

type commentWriter interface {
	WriteComment(ctx context.Context, article *models.Comment) error
}

type writeCommentService struct {
	repository commentWriter
}

func (s *writeCommentService) WriteComment(ctx context.Context, comment *models.Comment) error {
	if err := s.repository.WriteComment(ctx, comment); err != nil {
		return err
	}
	return nil
}
