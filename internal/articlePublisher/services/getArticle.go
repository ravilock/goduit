package services

import (
	"context"

	"github.com/ravilock/goduit/internal/articlePublisher/models"
)

type articleGetter interface {
	GetArticleBySlug(ctx context.Context, slug string) (*models.Article, error)
}

type getArticleService struct {
	repository articleGetter
}

func (s *getArticleService) GetArticleBySlug(ctx context.Context, slug string) (*models.Article, error) {
	article, err := s.repository.GetArticleBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	return article, nil
}
