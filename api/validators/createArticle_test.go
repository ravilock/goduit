package validators

import (
	"testing"

	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/api/requests"
)

func TestCreateArticle(t *testing.T) {
	t.Run("Title is required", func(t *testing.T) {
		request := assembleCreateArticleRequest()
		request.Article.Title = ""

		got := CreateArticle(request)
		assertError(t, got, api.RequiredFieldError("Title"))
	})
	t.Run("Title must not be blank", func(t *testing.T) {
		request := assembleCreateArticleRequest()
		request.Article.Title = "   "

		got := CreateArticle(request)
		assertError(t, got, api.RequiredFieldError("Title"))
	})
	t.Run("Title must have at least 5 chars", func(t *testing.T) {
		request := assembleCreateArticleRequest()
		request.Article.Title = "abcd"

		got := CreateArticle(request)
		assertError(t, got, api.InvalidFieldLength("Title", "min", "5"))
	})
	t.Run("Title must have at most 255 chars", func(t *testing.T) {
		request := assembleCreateArticleRequest()
		request.Article.Title = randomString(256)

		got := CreateArticle(request)
		assertError(t, got, api.InvalidFieldLength("Title", "max", "255"))
	})
	t.Run("Description is required", func(t *testing.T) {
		request := assembleCreateArticleRequest()
		request.Article.Description = ""

		got := CreateArticle(request)
		assertError(t, got, api.RequiredFieldError("Description"))
	})
	t.Run("Description must not be blank", func(t *testing.T) {
		request := assembleCreateArticleRequest()
		request.Article.Description = "   "

		got := CreateArticle(request)
		assertError(t, got, api.RequiredFieldError("Description"))
	})
	t.Run("Description must have at least 5 chars", func(t *testing.T) {
		request := assembleCreateArticleRequest()
		request.Article.Description = "abcd"

		got := CreateArticle(request)
		assertError(t, got, api.InvalidFieldLength("Description", "min", "5"))
	})
	t.Run("Description must have at most 255 chars", func(t *testing.T) {
		request := assembleCreateArticleRequest()
		request.Article.Description = randomString(256)

		got := CreateArticle(request)
		assertError(t, got, api.InvalidFieldLength("Description", "max", "255"))
	})
	t.Run("Body is required", func(t *testing.T) {
		request := assembleCreateArticleRequest()
		request.Article.Body = ""

		got := CreateArticle(request)
		assertError(t, got, api.RequiredFieldError("Body"))
	})
	t.Run("Body must not be blank", func(t *testing.T) {
		request := assembleCreateArticleRequest()
		request.Article.Body = "    "

		got := CreateArticle(request)
		assertError(t, got, api.RequiredFieldError("Body"))
	})
	t.Run("TagList must have unique tags", func(t *testing.T) {
		request := assembleCreateArticleRequest()
		request.Article.TagList = []string{"test", "test"}

		got := CreateArticle(request)
		assertError(t, got, api.UniqueFieldError("TagList"))
	})
	t.Run("Each tag in TagList must contain at least 3 chars", func(t *testing.T) {
		request := assembleCreateArticleRequest()
		request.Article.TagList = []string{"test", "tag", "f"}

		got := CreateArticle(request)
		assertError(t, got, api.InvalidFieldLength("TagList[2]", "min", "3"))
	})
	t.Run("TagList must contain at most 30 tags", func(t *testing.T) {
		request := assembleCreateArticleRequest()
		request.Article.TagList = append(request.Article.TagList, randomString(31))

		got := CreateArticle(request)
		assertError(t, got, api.InvalidFieldLength("TagList[3]", "max", "30"))
	})
}

func assembleCreateArticleRequest() *requests.CreateArticle {
	request := new(requests.CreateArticle)
	request.Article.Title = "Test Title"
	request.Article.Description = "Test Description"
	request.Article.Body = "Test Body"
	request.Article.TagList = []string{"test", "tag", "list"}
	return request
}
