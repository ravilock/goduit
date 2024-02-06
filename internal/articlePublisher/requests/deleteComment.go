package requests

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/ravilock/goduit/api/validators"
)

type DeleteCommentRequest struct {
	Slug string `param:"slug" validate:"required,notblank,min=5"`
	ID   string `param:"id" validate:"required,notblank"`
}

func (r *DeleteCommentRequest) Validate() error {
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
