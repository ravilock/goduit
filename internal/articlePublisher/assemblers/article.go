package assemblers

import (
	"github.com/ravilock/goduit/internal/articlePublisher/models"
	"github.com/ravilock/goduit/internal/articlePublisher/responses"
	profileManagerResponses "github.com/ravilock/goduit/internal/profileManager/responses"
)

func ArticleResponse(article *models.Article, author *profileManagerResponses.ProfileResponse) *responses.Article {
	response := new(responses.Article)
	response.Article.Slug = *article.Slug
	response.Article.Title = *article.Title
	response.Article.Description = *article.Description
	response.Article.Body = *article.Body
	response.Article.TagList = article.TagList
	response.Article.CreatedAt = *article.CreatedAt
	if article.UpdatedAt != nil {
		response.Article.UpdatedAt = *article.UpdatedAt
	}
	response.Article.Favorited = false
	response.Article.FavoritesCount = *article.FavoritesCount
	response.Article.Author = author.Profile
	return response
}
