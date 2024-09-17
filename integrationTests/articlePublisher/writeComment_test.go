package articlepublisher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	integrationtests "github.com/ravilock/goduit/integrationTests"
	articlePublisherRequests "github.com/ravilock/goduit/internal/articlePublisher/requests"
	articlePublisherResponses "github.com/ravilock/goduit/internal/articlePublisher/responses"
	profileManagerRequests "github.com/ravilock/goduit/internal/profileManager/requests"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

const (
	commentsPath = "comments"
)

func TestWriteComment(t *testing.T) {
	serverUrl := viper.GetString("server.url")
	writeCommentEndpoint := fmt.Sprintf("%s%s", serverUrl, "/api/articles")
	httpClient := http.Client{}

	t.Run("Should create a comment", func(t *testing.T) {
		// Arrange
		authorIdentity, authorToken := integrationtests.MustRegisterUser(t, profileManagerRequests.RegisterPayload{})
		article := integrationtests.MustWriteArticle(t, articlePublisherRequests.WriteArticlePayload{}, authorToken)
		writeCommentRequest := generateWriteCommentBody()
		requestBody, err := json.Marshal(writeCommentRequest)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s/%s", writeCommentEndpoint, article.Article.Slug, commentsPath), bytes.NewBuffer(requestBody))
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", authorToken))

		// Act
		res, err := httpClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, res.StatusCode)
		resBytes, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		writeCommentResponse := new(articlePublisherResponses.CommentResponse)
		err = json.Unmarshal(resBytes, writeCommentResponse)
		require.NoError(t, err)
		checkWriteCommentResponse(t, writeCommentRequest, authorIdentity.Username, writeCommentResponse)
	})

	t.Run("Should return HTTP 404 if targeted article does not exists", func(t *testing.T) {
		// Arrange
		_, authorToken := integrationtests.MustRegisterUser(t, profileManagerRequests.RegisterPayload{})
		writeCommentRequest := generateWriteCommentBody()
		requestBody, err := json.Marshal(writeCommentRequest)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s/%s", writeCommentEndpoint, uuid.NewString(), commentsPath), bytes.NewBuffer(requestBody))
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", authorToken))

		// Act
		res, err := httpClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusNotFound, res.StatusCode)
	})
}

func generateWriteCommentBody() *articlePublisherRequests.WriteCommentRequest {
	request := new(articlePublisherRequests.WriteCommentRequest)
	request.Comment.Body = "Test Comment Body"
	return request
}

func checkWriteCommentResponse(t *testing.T, request *articlePublisherRequests.WriteCommentRequest, authorUsername string, response *articlePublisherResponses.CommentResponse) {
	t.Helper()
	require.NotZero(t, response.Comment.ID)
	require.Equal(t, request.Comment.Body, response.Comment.Body, "Wrong comment body")
	require.Equal(t, authorUsername, response.Comment.Author.Username, "Wrong comment author username")
}
