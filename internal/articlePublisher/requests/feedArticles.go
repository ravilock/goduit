package requests

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/ravilock/goduit/api/validators"
)

type FeedArticlesRequest struct {
	Pagination FeedArticlesPagination
}

type FeedArticlesPagination struct {
	Limit  int `query:"limit" validate:"min=1,max=30"`
	Offset int `query:"offset" validate:"min=0"`
}

func NewFeedArticlesRequest() *FeedArticlesRequest {
	return &FeedArticlesRequest{
		FeedArticlesPagination{
			Limit: 20,
		},
	}
}

func (r *FeedArticlesRequest) Validate() error {
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
