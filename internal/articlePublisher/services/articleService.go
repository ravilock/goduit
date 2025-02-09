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
}

func NewArticlePublisher(articleRepository *repositories.ArticleRepository, articleQueuePublisher *publishers.ArticleQueuePublisher) *ArticlePublisher {
	write := writeArticleService{articleRepository, articleQueuePublisher}
	get := getArticleService{articleRepository}
	unpublish := unpublishArticleService{articleRepository}
	update := updateArticleService{articleRepository}
	list := listArticleService{articleRepository}
	return &ArticlePublisher{write, get, unpublish, update, list}
}
