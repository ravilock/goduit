package services

import (
	"context"

	"github.com/ravilock/goduit/internal/app/dtos"
	"github.com/ravilock/goduit/internal/app/repositories"
	"github.com/ravilock/goduit/internal/app/transformers"
)

func CreateArticle(article *dtos.Article, ctx context.Context) (*dtos.Article, error) {
	model := transformers.ArticleDtoToModel(article)

	if err := repositories.CreateArticle(model, ctx); err != nil {
		return nil, err
	}

	return transformers.ArticleModelToDto(model, article.Author), nil
}
