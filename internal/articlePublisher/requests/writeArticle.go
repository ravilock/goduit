package requests

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/ravilock/goduit/api/validators"
	"github.com/ravilock/goduit/internal/articlePublisher/models"
)

type WriteArticleRequest struct {
	Article WriteArticlePayload `json:"article" validate:"required"`
}

type WriteArticlePayload struct {
	Title       string   `json:"title" validate:"required,notblank,min=5,max=255"`
	Description string   `json:"description" validate:"required,notblank,min=5,max=255"`
	Body        string   `json:"body" validate:"required,notblank"`
	TagList     []string `json:"tagList" validate:"min=1,max=10,unique,dive,min=3,max=30"`
}

func (r *WriteArticleRequest) Model(authorID string) *models.Article {
	tags := deduplicateTags(r.Article.TagList)
	slug := makeSlug(r.Article.Title)
	return &models.Article{
		Author:         &authorID,
		Slug:           &slug,
		Title:          &r.Article.Title,
		Description:    &r.Article.Description,
		Body:           &r.Article.Body,
		TagList:        tags,
		FavoritesCount: new(int64),
	}
}

func deduplicateTags(tags []string) []string {
	tagMap := make(map[string]bool)
	deduplicatedTags := make([]string, 0, cap(tags))
	for _, tag := range tags {
		normalizedTag := strings.ToLower(tag)
		ok := tagMap[normalizedTag]
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
	return strings.ReplaceAll(loweredTitle, " ", "-")
}

func (r *WriteArticleRequest) Validate() error {
	if err := validators.Validate.Struct(r); err != nil {
		if validationErrors := new(validator.ValidationErrors); errors.As(err, validationErrors) {
			for _, validationError := range *validationErrors {
				return validators.ToHTTP(validationError)
			}
		}
		return err
	}
	return nil
}
