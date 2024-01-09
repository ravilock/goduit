package requests

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/ravilock/goduit/api/validators"
)

type ListArticlesRequest struct {
	Pagination ListArticlesPagination
	Filters    ListArticlesFilters
}

type ListArticlesFilters struct {
	Tag       string `query:"tag" validate:"omitempty,notblank,min=3,max=30"`
	Author    string `query:"author" validate:"omitempty,notblank,min=5,max=255"`
	Favorited string
}

type ListArticlesPagination struct {
	Limit  int `query:"limit" validate:"min=1,max=30"`
	Offset int `query:"offset" validate:"min=0"`
}

func NewListArticlesRequest() *ListArticlesRequest {
	return &ListArticlesRequest{
		ListArticlesPagination{
			Limit: 20,
		},
		ListArticlesFilters{},
	}
}

func (r *ListArticlesRequest) Validate() error {
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
