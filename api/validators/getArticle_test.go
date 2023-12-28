package validators

import (
	"testing"

	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/api/requests"
)

func TestGetArticle(t *testing.T) {
	InitValidator()
	t.Run("Slug is required", func(t *testing.T) {
		request := &requests.GetArticle{}

		got := GetArticle(request)
		assertError(t, got, api.RequiredFieldError("Slug"))
	})
	t.Run("Slug must not be blank", func(t *testing.T) {
		request := &requests.GetArticle{Slug: "   "}

		got := GetArticle(request)
		assertError(t, got, api.RequiredFieldError("Slug"))
	})
	t.Run("Slug must have at least 5 chars", func(t *testing.T) {
		request := &requests.GetArticle{Slug: "1234"}

		got := GetArticle(request)
		assertError(t, got, api.InvalidFieldLength("Slug", "min", "5"))
	})
	t.Run("Slug must have at most 255 chars", func(t *testing.T) {
		request := &requests.GetArticle{Slug: randomString(256)}

		got := GetArticle(request)
		assertError(t, got, api.InvalidFieldLength("Slug", "max", "255"))
	})
}
