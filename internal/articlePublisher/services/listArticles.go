package services

import (
	"context"

	"github.com/ravilock/goduit/internal/articlePublisher/models"
)

type articleLister interface {
	ListArticles(ctx context.Context, author, tag string, limit, offset int64) ([]*models.Article, error)
}

type ListArticlesService struct {
	repository articleLister
}

func NewListArticlesService(repository articleLister) *ListArticlesService {
	return &ListArticlesService{
		repository: repository,
	}
}

func (s *ListArticlesService) ListArticles(ctx context.Context, author, tag string, limit, offset int64) ([]*models.Article, error) {
	return s.repository.ListArticles(ctx, author, tag, limit, offset)
}
