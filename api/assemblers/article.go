package assemblers

import (
	"strings"
	"time"

	"github.com/ravilock/goduit/api/requests"
	"github.com/ravilock/goduit/api/responses"
	"github.com/ravilock/goduit/internal/app/dtos"
)

func CreateArticle(request *requests.CreateArticle, author *dtos.Profile) *dtos.Article {
	tags := deduplicateTags(request.Article.TagList)
	slug := makeSlug(request.Article.Title)
	createAt := time.Now()
	return &dtos.Article{
		Slug:           &slug,
		Title:          &request.Article.Title,
		Description:    &request.Article.Description,
		Body:           &request.Article.Body,
		TagList:        &tags,
		CreatedAt:      &createAt,
		UpdatedAt:      nil,
		Favorited:      false,
		FavoritesCount: 0,
		Author:         author,
	}
}

func deduplicateTags(tags []string) []string {
	tagMap := make(map[string]bool)
	deduplicatedTags := make([]string, 0, cap(tags))
	for _, tag := range tags {
		normalizedTag := strings.ToLower(tag)
		ok, _ := tagMap[normalizedTag]
		if ok {
			continue
		}
		tagMap[normalizedTag] = true
		deduplicatedTags = append(deduplicatedTags, normalizedTag)
	}
	return deduplicatedTags
}

func makeSlug(title string) string {
	loweredTitle := strings.ToLower(title)
	titleWords := strings.Split(loweredTitle, " ")
	return strings.Join(titleWords, "-")
}

func ArticleResponse(dto *dtos.Article) *responses.Article {
	var author responses.Profile
	author.Username = *dto.Author.Username
	if dto.Author.Bio != nil {
		author.Bio = *dto.Author.Bio
	}
	if dto.Author.Image != nil {
		author.Image = *dto.Author.Image
	}
	response := new(responses.Article)
	response.Article.Slug = *dto.Slug
	response.Article.Title = *dto.Title
	response.Article.Description = *dto.Description
	response.Article.Body = *dto.Body
	response.Article.TagList = *dto.TagList
	response.Article.CreatedAt = *dto.CreatedAt
	if dto.UpdatedAt != nil {
		response.Article.UpdatedAt = *dto.UpdatedAt
	}
	response.Article.Favorited = dto.Favorited
	response.Article.FavoritesCount = dto.FavoritesCount
	response.Article.Author = author
	return response
}
