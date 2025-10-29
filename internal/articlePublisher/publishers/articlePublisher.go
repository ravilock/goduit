package publishers

import (
	"context"

	"github.com/ravilock/goduit/internal/articlePublisher/models"
	"github.com/ravilock/goduit/internal/queue"
)

type ArticlePublisher struct {
	publisher queue.Publisher
}

func NewArticlePublisher(publisher queue.Publisher) *ArticlePublisher {
	return &ArticlePublisher{
		publisher: publisher,
	}
}

func (p *ArticlePublisher) PublishArticle(ctx context.Context, article *models.Article) error {
	return p.publisher.Publish(ctx, []byte(article.ID.Hex()))
}
