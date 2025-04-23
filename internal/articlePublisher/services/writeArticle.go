package services

import (
	"context"

	"github.com/ravilock/goduit/internal/articlePublisher/models"
)

type articleWriter interface {
	WriteArticle(ctx context.Context, article *models.Article) error
}

type articlePublisher interface {
	PublishArticle(ctx context.Context, article *models.Article) error
}

type WriteArticleService struct {
	repository articleWriter
	queue      articlePublisher
}

func NewWriteArticleService(repository articleWriter, queue articlePublisher) *WriteArticleService {
	return &WriteArticleService{
		repository: repository,
		queue:      queue,
	}
}

func (s *WriteArticleService) WriteArticle(ctx context.Context, article *models.Article) error {
	if err := s.repository.WriteArticle(ctx, article); err != nil {
		return err
	}
	return s.queue.PublishArticle(ctx, article)
}
