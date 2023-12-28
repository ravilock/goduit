package requests

import (
	"testing"

	"github.com/ravilock/goduit/api"
	"github.com/stretchr/testify/assert"
)

func TestArticleSlug(t *testing.T) {
	t.Run("Valid request should not return errors", func(t *testing.T) {
		request := generateArticleSlugRequest()
		err := request.Validate()
		assert.NoError(t, err)
	})
	t.Run("Slug is required", func(t *testing.T) {
		request := generateArticleSlugRequest()
		request.Slug = ""
		err := request.Validate()
		assert.ErrorContains(t, err, api.RequiredFieldError("Slug").Error())
	})
	t.Run("Slug should not be blank", func(t *testing.T) {
		request := generateArticleSlugRequest()
		request.Slug = " "
		err := request.Validate()
		assert.ErrorContains(t, err, api.RequiredFieldError("Slug").Error())
	})
	t.Run("Slug should contain at least 5 chars", func(t *testing.T) {
		request := generateArticleSlugRequest()
		request.Slug = "1234"
		err := request.Validate()
		assert.ErrorContains(t, err, api.InvalidFieldLength("Slug", "min", "5").Error())
	})
	t.Run("Slug should contain at most 255 chars", func(t *testing.T) {
		request := generateArticleSlugRequest()
		request.Slug = randomString(256)
		err := request.Validate()
		assert.ErrorContains(t, err, api.InvalidFieldLength("Slug", "max", "255").Error())
	})
}

func generateArticleSlugRequest() *ArticleSlug {
	return &ArticleSlug{
		Slug: "test-slug",
	}
}
