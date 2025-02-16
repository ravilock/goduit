package services

import (
	"github.com/ravilock/goduit/internal/articlePublisher/publishers"
	"github.com/ravilock/goduit/internal/articlePublisher/repositories"
)

type ArticlePublisher struct {
	writeArticleService
	getArticleService
	unpublishArticleService
	updateArticleService
	listArticleService
	feedArticlesService
}

func NewArticlePublisher(articleRepository *repositories.ArticleRepository, feedRepository *repositories.FeedRepository, articleQueuePublisher *publishers.ArticleQueuePublisher) *ArticlePublisher {
	write := writeArticleService{articleRepository, articleQueuePublisher}
	get := getArticleService{articleRepository}
	unpublish := unpublishArticleService{articleRepository}
	update := updateArticleService{articleRepository}
	list := listArticleService{articleRepository}
	feed := feedArticlesService{articleRepository, feedRepository}
	return &ArticlePublisher{write, get, unpublish, update, list, feed}
}
