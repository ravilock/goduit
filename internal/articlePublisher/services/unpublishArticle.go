package services

import "context"

type articleDeleter interface {
	DeleteArticle(ctx context.Context, slug string) error
}

type UnpublishArticleService struct {
	repository articleDeleter
}

func NewUnpublishArticleService(repository articleDeleter) *UnpublishArticleService {
	return &UnpublishArticleService{
		repository: repository,
	}
}

func (s *UnpublishArticleService) UnpublishArticle(ctx context.Context, slug string) error {
	return s.repository.DeleteArticle(ctx, slug)
}
