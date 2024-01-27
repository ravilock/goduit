package requests

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/ravilock/goduit/api/validators"
	"github.com/ravilock/goduit/internal/articlePublisher/models"
)

type WriteCommentRequest struct {
	Slug    string              `param:"slug" validate:"required,notblank,min=5"`
	Comment WriteCommentPayload `json:"comment" validate:"required"`
}

type WriteCommentPayload struct {
	Body string `json:"body" validate:"required,notblank,min=5,max=255"`
}

func (r *WriteCommentRequest) Model(authorID string) *models.Comment {
	createdAt := time.Now()
	return &models.Comment{
		Author:    &authorID,
		Body:      &r.Comment.Body,
		CreatedAt: &createdAt,
		UpdatedAt: nil,
	}
}

func (r *WriteCommentRequest) Validate() error {
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
