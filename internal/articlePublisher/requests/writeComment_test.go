package requests

import (
	"testing"

	"github.com/ravilock/goduit/api"
	"github.com/stretchr/testify/require"
)

func TestWriteComment(t *testing.T) {
	t.Run("Valid requests should not return errors", func(t *testing.T) {
		request := generateWriteCommentRequest()
		err := request.Validate()
		require.NoError(t, err)
	})
	t.Run("Body is required", func(t *testing.T) {
		request := generateWriteCommentRequest()
		request.Comment.Body = ""
		err := request.Validate()
		require.ErrorContains(t, err, api.RequiredFieldError("Body").Error())
	})
	t.Run("Body should not be blank", func(t *testing.T) {
		request := generateWriteCommentRequest()
		request.Comment.Body = " "
		err := request.Validate()
		require.ErrorContains(t, err, api.RequiredFieldError("Body").Error())
	})
	t.Run("Body should contain at least 5 chars", func(t *testing.T) {
		request := generateWriteCommentRequest()
		request.Comment.Body = "1234"
		err := request.Validate()
		require.ErrorContains(t, err, api.InvalidFieldLimit("Body", "min", "5").Error())
	})
	t.Run("Body should contain at least 140 chars", func(t *testing.T) {
		request := generateWriteCommentRequest()
		request.Comment.Body = randomString(141)
		err := request.Validate()
		require.ErrorContains(t, err, api.InvalidFieldLimit("Body", "max", "140").Error())
	})
	t.Run("Slug is required", func(t *testing.T) {
		request := generateWriteCommentRequest()
		request.Slug = ""
		err := request.Validate()
		require.ErrorContains(t, err, api.RequiredFieldError("Slug").Error())
	})
	t.Run("Slug should not be blank", func(t *testing.T) {
		request := generateWriteCommentRequest()
		request.Slug = " "
		err := request.Validate()
		require.ErrorContains(t, err, api.RequiredFieldError("Slug").Error())
	})
}

func generateWriteCommentRequest() *WriteCommentRequest {
	comment := new(WriteCommentRequest)
	comment.Slug = "test-article-slug"
	comment.Comment.Body = "Test Body"
	return comment
}
