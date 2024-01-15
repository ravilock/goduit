package services

import (
	"context"

	"github.com/ravilock/goduit/internal/articlePublisher/models"
)

type articleLister interface {
	ListArticles(ctx context.Context, author, tag string, limit, offset int64) ([]*models.Article, error)
}

type listArticleService struct {
	repository articleLister
}

func (s *listArticleService) ListArticles(ctx context.Context, author, tag string, limit, offset int64) ([]*models.Article, error) {
	return s.repository.ListArticles(ctx, author, tag, limit, offset)
}
