package requests

import (
	"testing"

	"github.com/ravilock/goduit/api"
	"github.com/stretchr/testify/require"
)

func TestArticleSlug(t *testing.T) {
	t.Run("Valid request should not return errors", func(t *testing.T) {
		request := generateArticleSlugRequest()
		err := request.Validate()
		require.NoError(t, err)
	})
	t.Run("Slug is required", func(t *testing.T) {
		request := generateArticleSlugRequest()
		request.Slug = ""
		err := request.Validate()
		require.ErrorContains(t, err, api.RequiredFieldError("Slug").Error())
	})
	t.Run("Slug should not be blank", func(t *testing.T) {
		request := generateArticleSlugRequest()
		request.Slug = " "
		err := request.Validate()
		require.ErrorContains(t, err, api.RequiredFieldError("Slug").Error())
	})
	t.Run("Slug should contain at least 5 chars", func(t *testing.T) {
		request := generateArticleSlugRequest()
		request.Slug = "1234"
		err := request.Validate()
		require.ErrorContains(t, err, api.InvalidFieldLimit("Slug", "min", "5").Error())
	})
}

func generateArticleSlugRequest() *ArticleSlugRequest {
	return &ArticleSlugRequest{
		Slug: "test-slug",
	}
}
