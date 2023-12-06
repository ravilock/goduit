package services

import "context"

type articleDeleter interface {
	DeleteArticle(ctx context.Context, slug string) error
}

type unpublishArticleService struct {
	repository articleDeleter
}

func (s *unpublishArticleService) UnpublishArticle(ctx context.Context, slug string) error {
	return s.repository.DeleteArticle(ctx, slug)
}
