package requests

import (
	"testing"

	"github.com/ravilock/goduit/api"
	"github.com/stretchr/testify/assert"
)

func TestUpdateArticle(t *testing.T) {
	t.Run("Valid request should not return errors", func(t *testing.T) {
		request := generateUpdateArticleRequest()
		err := request.Validate()
		assert.NoError(t, err)
	})
	t.Run("Slug is required", func(t *testing.T) {
		request := generateUpdateArticleRequest()
		request.Article.Slug = ""
		err := request.Validate()
		assert.ErrorContains(t, err, api.RequiredFieldError("Slug").Error())
	})
	t.Run("Slug should not be blank", func(t *testing.T) {
		request := generateUpdateArticleRequest()
		request.Article.Slug = " "
		err := request.Validate()
		assert.ErrorContains(t, err, api.RequiredFieldError("Slug").Error())
	})
	t.Run("Slug should contain at least 5 chars", func(t *testing.T) {
		request := generateUpdateArticleRequest()
		request.Article.Slug = "1234"
		err := request.Validate()
		assert.ErrorContains(t, err, api.InvalidFieldLength("Slug", "min", "5").Error())
	})
	t.Run("Title is required", func(t *testing.T) {
		request := generateUpdateArticleRequest()
		request.Article.Title = ""
		err := request.Validate()
		assert.ErrorContains(t, err, api.RequiredFieldError("Title").Error())
	})
	t.Run("Title should not be blank", func(t *testing.T) {
		request := generateUpdateArticleRequest()
		request.Article.Title = " "
		err := request.Validate()
		assert.ErrorContains(t, err, api.RequiredFieldError("Title").Error())
	})
	t.Run("Title should contain at least 5 chars", func(t *testing.T) {
		request := generateUpdateArticleRequest()
		request.Article.Title = "1234"
		err := request.Validate()
		assert.ErrorContains(t, err, api.InvalidFieldLength("Title", "min", "5").Error())
	})
	t.Run("Title should contain at most 255 chars", func(t *testing.T) {
		request := generateUpdateArticleRequest()
		request.Article.Title = randomString(256)
		err := request.Validate()
		assert.ErrorContains(t, err, api.InvalidFieldLength("Title", "max", "255").Error())
	})
	t.Run("Description is required", func(t *testing.T) {
		request := generateUpdateArticleRequest()
		request.Article.Description = ""
		err := request.Validate()
		assert.ErrorContains(t, err, api.RequiredFieldError("Description").Error())
	})
	t.Run("Description should not be blank", func(t *testing.T) {
		request := generateUpdateArticleRequest()
		request.Article.Description = " "
		err := request.Validate()
		assert.ErrorContains(t, err, api.RequiredFieldError("Description").Error())
	})
	t.Run("Description should contain at least 5 chars", func(t *testing.T) {
		request := generateUpdateArticleRequest()
		request.Article.Description = "1234"
		err := request.Validate()
		assert.ErrorContains(t, err, api.InvalidFieldLength("Description", "min", "5").Error())
	})
	t.Run("Description should contain at most 255 chars", func(t *testing.T) {
		request := generateUpdateArticleRequest()
		request.Article.Description = randomString(256)
		err := request.Validate()
		assert.ErrorContains(t, err, api.InvalidFieldLength("Description", "max", "255").Error())
	})
	t.Run("Body is required", func(t *testing.T) {
		request := generateUpdateArticleRequest()
		request.Article.Body = ""
		err := request.Validate()
		assert.ErrorContains(t, err, api.RequiredFieldError("Body").Error())
	})
	t.Run("Body should not be blank", func(t *testing.T) {
		request := generateUpdateArticleRequest()
		request.Article.Body = " "
		err := request.Validate()
		assert.ErrorContains(t, err, api.RequiredFieldError("Body").Error())
	})
}

func generateUpdateArticleRequest() *UpdateArticle {
	article := new(UpdateArticle)
	article.Article.Slug = "test-slug"
	article.Article.Title = "Test Title"
	article.Article.Description = "Test Description"
	article.Article.Body = "Test Body"
	return article
}
