package services

import (
	"context"

	"github.com/ravilock/goduit/internal/articlePublisher/models"
)

type articleUpdater interface {
	UpdateArticle(ctx context.Context, slug string, article *models.Article) error
}

type UpdateArticleService struct {
	repository articleUpdater
}

func NewUpdateArticleService(repository articleUpdater) *UpdateArticleService {
	return &UpdateArticleService{
		repository: repository,
	}
}

func (s *UpdateArticleService) UpdateArticle(ctx context.Context, slug string, article *models.Article) error {
	return s.repository.UpdateArticle(ctx, slug, article)
}
