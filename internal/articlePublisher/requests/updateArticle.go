package requests

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/ravilock/goduit/api/validators"
	"github.com/ravilock/goduit/internal/articlePublisher/models"
)

type UpdateArticle struct {
	Slug    string               `param:"slug" validate:"required,notblank,min=5"`
	Article UpdateArticlePayload `json:"article" validate:"required"`
}

type UpdateArticlePayload struct {
	Title       string `json:"title" validate:"required,notblank,min=5,max=255"`
	Description string `json:"description" validate:"required,notblank,min=5,max=255"`
	Body        string `json:"body" validate:"required,notblank"`
}

func (r *UpdateArticle) Model() *models.Article {
	slug := makeSlug(r.Article.Title)
	updatedAt := time.Now()
	return &models.Article{
		Author:         nil,
		Slug:           &slug,
		Title:          &r.Article.Title,
		Description:    &r.Article.Description,
		Body:           &r.Article.Body,
		TagList:        nil,
		CreatedAt:      nil,
		UpdatedAt:      &updatedAt,
		FavoritesCount: nil,
	}
}

func (r *UpdateArticle) Validate() error {
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
