package services

import (
	"context"

	"github.com/ravilock/goduit/internal/articlePublisher/models"
)

type articleWriter interface {
	WriteArticle(ctx context.Context, article *models.Article) error
}

type writeArticleService struct {
	repository articleWriter
}

func (s *writeArticleService) WriteArticle(ctx context.Context, article *models.Article) error {
	if err := s.repository.WriteArticle(ctx, article); err != nil {
		return err
	}
	return nil
}
