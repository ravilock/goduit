package requests

import (
	"errors"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/ravilock/goduit/api/validators"
	"github.com/ravilock/goduit/internal/articlePublisher/models"
)

type WriteArticle struct {
	Article struct {
		Title       string   `json:"title" validate:"required,notblank,min=5,max=255"`
		Description string   `json:"description" validate:"required,notblank,min=5,max=255"`
		Body        string   `json:"body" validate:"required,notblank"`
		TagList     []string `json:"tagList" validate:"min=1,max=10,unique,dive,min=3,max=30"`
	} `json:"article" validate:"required"`
}

func (r *WriteArticle) Model(author string) *models.Article {
	tags := deduplicateTags(r.Article.TagList)
	slug := makeSlug(r.Article.Title)
	createAt := time.Now()
	return &models.Article{
		Author:         &author,
		Slug:           &slug,
		Title:          &r.Article.Title,
		Description:    &r.Article.Description,
		Body:           &r.Article.Body,
		TagList:        &tags,
		CreatedAt:      &createAt,
		UpdatedAt:      nil,
		FavoritesCount: 0,
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
	titleWords := strings.Split(loweredTitle, " ")
	return strings.Join(titleWords, "-")
}

func (r *WriteArticle) Validate() error {
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
