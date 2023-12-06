package services

import (
	"context"

	"github.com/ravilock/goduit/internal/articlePublisher/models"
)

type articleUpdater interface {
	UpdateArticle(ctx context.Context, slug string, article *models.Article) error
}

type updateArticleService struct {
	repository articleUpdater
}

func (s *updateArticleService) UpdateArticle(ctx context.Context, slug string, article *models.Article) error {
	return s.repository.UpdateArticle(ctx, slug, article)
}
