package requests

import (
	"testing"

	"github.com/ravilock/goduit/api"
	"github.com/stretchr/testify/require"
)

func TestListArticles(t *testing.T) {
	t.Run("Valid request should return errors", func(t *testing.T) {
		request := generateListArticlesRequest()
		err := request.Validate()
		require.NoError(t, err)
	})
	t.Run("Tag is optional, but should not be blank", func(t *testing.T) {
		request := generateListArticlesRequest()
		request.Filters.Tag = ""
		err := request.Validate()
		require.NoError(t, err)
		request.Filters.Tag = " "
		err = request.Validate()
		require.ErrorContains(t, err, api.RequiredFieldError("Tag").Error())
	})
	t.Run("Tag should contain at least 3 chars", func(t *testing.T) {
		request := generateListArticlesRequest()
		request.Filters.Tag = "12"
		err := request.Validate()
		require.ErrorContains(t, err, api.InvalidFieldLimit("Tag", "min", "3").Error())
	})
	t.Run("Tag should contain at most 30 chars", func(t *testing.T) {
		request := generateListArticlesRequest()
		request.Filters.Tag = randomString(31)
		err := request.Validate()
		require.ErrorContains(t, err, api.InvalidFieldLimit("Tag", "max", "30").Error())
	})
	t.Run("Author is optional, but should not be blank", func(t *testing.T) {
		request := generateListArticlesRequest()
		request.Filters.Author = ""
		err := request.Validate()
		require.NoError(t, err)
		request.Filters.Author = " "
		err = request.Validate()
		require.ErrorContains(t, err, api.RequiredFieldError("Author").Error())
	})
	t.Run("Author should contain at least 5 chars", func(t *testing.T) {
		request := generateListArticlesRequest()
		request.Filters.Author = "1234"
		err := request.Validate()
		require.ErrorContains(t, err, api.InvalidFieldLimit("Author", "min", "5").Error())
	})
	t.Run("Author should contain at most 255 chars", func(t *testing.T) {
		request := generateListArticlesRequest()
		request.Filters.Author = randomString(256)
		err := request.Validate()
		require.ErrorContains(t, err, api.InvalidFieldLimit("Author", "max", "255").Error())
	})
	t.Run("Limit should have min value 1", func(t *testing.T) {
		request := generateListArticlesRequest()
		request.Pagination.Limit = 0
		err := request.Validate()
		require.ErrorContains(t, err, api.InvalidFieldLimit("Limit", "min", "1").Error())
	})
	t.Run("Limit should have max value 30", func(t *testing.T) {
		request := generateListArticlesRequest()
		request.Pagination.Limit = 31
		err := request.Validate()
		require.ErrorContains(t, err, api.InvalidFieldLimit("Limit", "max", "30").Error())
	})
	t.Run("Offset should have min value 1", func(t *testing.T) {
		request := generateListArticlesRequest()
		request.Pagination.Offset = -1
		err := request.Validate()
		require.ErrorContains(t, err, api.InvalidFieldLimit("Offset", "min", "0").Error())
	})
}

func generateListArticlesRequest() *ListArticlesRequest {
	return &ListArticlesRequest{
		ListArticlesPagination{
			Limit:  20,
			Offset: 20,
		},
		ListArticlesFilters{
			Tag:    "test-tag",
			Author: "test-ahutor",
		},
	}
}
