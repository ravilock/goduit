package services

import (
	"context"

	"github.com/ravilock/goduit/internal/articlePublisher/models"
)

type articlesGetter interface {
	GetArticlesByIDs(ctx context.Context, IDs []string) ([]*models.Article, error)
}

type feedPaginator interface {
	PaginateFeed(ctx context.Context, user string, limit, offset int64) ([]models.FeedFragment, error)
}

type FeedArticlesService struct {
	repository     articlesGetter
	feedRepository feedPaginator
}

func (s *FeedArticlesService) FeedArticles(ctx context.Context, user string, limit, offset int64) ([]*models.Article, error) {
	feedFragments, err := s.feedRepository.PaginateFeed(ctx, user, limit, offset)
	if err != nil {
		return nil, err
	}
	articleIDs := make([]string, len(feedFragments))
	for i, fragment := range feedFragments {
		articleIDs[i] = *fragment.ArticleID
	}
	return s.repository.GetArticlesByIDs(ctx, articleIDs)
}
