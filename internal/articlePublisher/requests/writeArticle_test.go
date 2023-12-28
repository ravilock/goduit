package requests

import (
	"testing"

	"github.com/ravilock/goduit/api"
	"github.com/stretchr/testify/assert"
)

func TestWriteArticle(t *testing.T) {
	t.Run("Valid request should not return errors", func(t *testing.T) {
		request := generateWriteArticleRequest()
		err := request.Validate()
		assert.NoError(t, err)
	})
	t.Run("Title is required", func(t *testing.T) {
		request := generateWriteArticleRequest()
		request.Article.Title = ""
		err := request.Validate()
		assert.ErrorContains(t, err, api.RequiredFieldError("Title").Error())
	})
	t.Run("Title should not be blank", func(t *testing.T) {
		request := generateWriteArticleRequest()
		request.Article.Title = " "
		err := request.Validate()
		assert.ErrorContains(t, err, api.RequiredFieldError("Title").Error())
	})
	t.Run("Title should contain at least 5 chars", func(t *testing.T) {
		request := generateWriteArticleRequest()
		request.Article.Title = "1234"
		err := request.Validate()
		assert.ErrorContains(t, err, api.InvalidFieldLength("Title", "min", "5").Error())
	})
	t.Run("Title should contain at most 255 chars", func(t *testing.T) {
		request := generateWriteArticleRequest()
		request.Article.Title = randomString(256)
		err := request.Validate()
		assert.ErrorContains(t, err, api.InvalidFieldLength("Title", "max", "255").Error())
	})
	t.Run("Description is required", func(t *testing.T) {
		request := generateWriteArticleRequest()
		request.Article.Description = ""
		err := request.Validate()
		assert.ErrorContains(t, err, api.RequiredFieldError("Description").Error())
	})
	t.Run("Description should not be blank", func(t *testing.T) {
		request := generateWriteArticleRequest()
		request.Article.Description = " "
		err := request.Validate()
		assert.ErrorContains(t, err, api.RequiredFieldError("Description").Error())
	})
	t.Run("Description should contain at least 5 chars", func(t *testing.T) {
		request := generateWriteArticleRequest()
		request.Article.Description = "1234"
		err := request.Validate()
		assert.ErrorContains(t, err, api.InvalidFieldLength("Description", "min", "5").Error())
	})
	t.Run("Description should contain at most 255 chars", func(t *testing.T) {
		request := generateWriteArticleRequest()
		request.Article.Description = randomString(256)
		err := request.Validate()
		assert.ErrorContains(t, err, api.InvalidFieldLength("Description", "max", "255").Error())
	})
	t.Run("Body is required", func(t *testing.T) {
		request := generateWriteArticleRequest()
		request.Article.Body = ""
		err := request.Validate()
		assert.ErrorContains(t, err, api.RequiredFieldError("Body").Error())
	})
	t.Run("Body should not be blank", func(t *testing.T) {
		request := generateWriteArticleRequest()
		request.Article.Body = " "
		err := request.Validate()
		assert.ErrorContains(t, err, api.RequiredFieldError("Body").Error())
	})
	t.Run("TagList should be unique", func(t *testing.T) {
		request := generateWriteArticleRequest()
		request.Article.TagList = []string{"test tag", "test tag"}
		err := request.Validate()
		assert.ErrorContains(t, err, api.UniqueFieldError("TagList").Error())
	})
	t.Run("Should require at least one tag in TagList", func(t *testing.T) {
		request := generateWriteArticleRequest()
		request.Article.TagList = []string{}
		err := request.Validate()
		assert.ErrorContains(t, err, api.InvalidFieldLength("TagList", "min", "1").Error())
	})
	t.Run("TagList should have at most 10 tags", func(t *testing.T) {
		request := generateWriteArticleRequest()
		request.Article.TagList = []string{}
		for i := 0; i < 11; i++ {
			request.Article.TagList = append(request.Article.TagList, randomString(30))
		}
		err := request.Validate()
		assert.ErrorContains(t, err, api.InvalidFieldLength("TagList", "max", "10").Error())
	})
	t.Run("Each Tag on TagList should have at least 3 chars", func(t *testing.T) {
		request := generateWriteArticleRequest()
		request.Article.TagList = []string{"12"}
		err := request.Validate()
		assert.ErrorContains(t, err, api.InvalidFieldLength("TagList[0]", "min", "3").Error())
	})
	t.Run("Each Tag on TagList should have at most 30 chars", func(t *testing.T) {
		request := generateWriteArticleRequest()
		request.Article.TagList = []string{randomString(31)}
		err := request.Validate()
		assert.ErrorContains(t, err, api.InvalidFieldLength("TagList[0]", "max", "30").Error())
	})
}

func generateWriteArticleRequest() *WriteArticle {
	article := new(WriteArticle)
	article.Article.Title = "Test Title"
	article.Article.Description = "Test Description"
	article.Article.Body = "Test Body"
	article.Article.TagList = []string{"Test Tag"}
	return article
}
