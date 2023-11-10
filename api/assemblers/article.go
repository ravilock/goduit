package assemblers

import (
	"github.com/ravilock/goduit/api/responses"
	"github.com/ravilock/goduit/internal/app/models"
)

func ArticleResponse(article *models.Article, author responses.ProfileResponse) *responses.Article {
	response := new(responses.Article)
	response.Article.Slug = *article.Slug
	response.Article.Title = *article.Title
	response.Article.Description = *article.Description
	response.Article.Body = *article.Body
	response.Article.TagList = *article.TagList
	response.Article.CreatedAt = *article.CreatedAt
	if article.UpdatedAt != nil {
		response.Article.UpdatedAt = *article.UpdatedAt
	}
	response.Article.Favorited = false
	response.Article.FavoritesCount = article.FavoritesCount
	response.Article.Author = author.Profile
	return response
}
