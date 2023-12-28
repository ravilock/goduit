package services

import "github.com/ravilock/goduit/internal/articlePublisher/repositories"

type ArticlePublisher struct {
	writeArticleService
}

func NewArticlePublisher(articleRepository *repositories.ArticleRepository) *ArticlePublisher {
	write := writeArticleService{articleRepository}
	return &ArticlePublisher{write}
}
