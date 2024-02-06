package requests

import (
	"testing"

	"github.com/google/uuid"
	"github.com/ravilock/goduit/api"
	"github.com/stretchr/testify/require"
)

func TestDeleteComment(t *testing.T) {
	t.Run("Valid request should not return errors", func(t *testing.T) {
		request := generateDeleteCommentRequest()
		err := request.Validate()
		require.NoError(t, err)
	})
	t.Run("Slug is required", func(t *testing.T) {
		request := generateDeleteCommentRequest()
		request.Slug = ""
		err := request.Validate()
		require.ErrorContains(t, err, api.RequiredFieldError("Slug").Error())
	})
	t.Run("Slug should not be blank", func(t *testing.T) {
		request := generateDeleteCommentRequest()
		request.Slug = " "
		err := request.Validate()
		require.ErrorContains(t, err, api.RequiredFieldError("Slug").Error())
	})
	t.Run("Slug should contain at least 5 chars", func(t *testing.T) {
		request := generateDeleteCommentRequest()
		request.Slug = "1234"
		err := request.Validate()
		require.ErrorContains(t, err, api.InvalidFieldLimit("Slug", "min", "5").Error())
	})
	t.Run("ID is required", func(t *testing.T) {
		request := generateDeleteCommentRequest()
		request.ID = ""
		err := request.Validate()
		require.ErrorContains(t, err, api.RequiredFieldError("ID").Error())
	})
	t.Run("ID should not be blank", func(t *testing.T) {
		request := generateDeleteCommentRequest()
		request.ID = " "
		err := request.Validate()
		require.ErrorContains(t, err, api.RequiredFieldError("ID").Error())
	})
}

func generateDeleteCommentRequest() *DeleteCommentRequest {
	return &DeleteCommentRequest{
		Slug: "test-slug",
		ID:   uuid.NewString(),
	}
}
