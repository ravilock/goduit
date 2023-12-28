package transformers

import (
	"github.com/ravilock/goduit/internal/app/dtos"
	"github.com/ravilock/goduit/internal/app/models"
)

func ArticleDtoToModel(article *dtos.Article) *models.Article {
	return &models.Article{
		Author:         article.Author.Username,
		Slug:           article.Slug,
		Title:          article.Title,
		Description:    article.Description,
		Body:           article.Body,
		TagList:        article.TagList,
		CreatedAt:      article.CreatedAt,
		UpdatedAt:      article.UpdatedAt,
		FavoritesCount: article.FavoritesCount,
	}
}

func ArticleModelToDto(article *models.Article, profile *dtos.Profile) *dtos.Article {
	return &dtos.Article{
		Slug:           article.Slug,
		Title:          article.Title,
		Description:    article.Description,
		Body:           article.Body,
		TagList:        article.TagList,
		CreatedAt:      article.CreatedAt,
		UpdatedAt:      article.UpdatedAt,
		FavoritesCount: article.FavoritesCount,
		Author:         profile,
	}
}
