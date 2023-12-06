package services

import "github.com/ravilock/goduit/internal/articlePublisher/repositories"

type ArticlePublisher struct {
	writeArticleService
	getArticleService
	unpublishArticleService
	updateArticleService
}

func NewArticlePublisher(articleRepository *repositories.ArticleRepository) *ArticlePublisher {
	write := writeArticleService{articleRepository}
	get := getArticleService{articleRepository}
	unpublish := unpublishArticleService{articleRepository}
	updated := updateArticleService{articleRepository}
	return &ArticlePublisher{write, get, unpublish, updated}
}
