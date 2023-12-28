package services

import "github.com/ravilock/goduit/internal/articlePublisher/repositories"

type ArticlePublisher struct {
	writeArticleService
	getArticleService
}

func NewArticlePublisher(articleRepository *repositories.ArticleRepository) *ArticlePublisher {
	write := writeArticleService{articleRepository}
	get := getArticleService{articleRepository}
	return &ArticlePublisher{write, get}
}
