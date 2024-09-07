package articlepublisher

import (
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

func TestListComments(t *testing.T) {
	serverUrl := viper.GetString("server.url")
	listCommentsEndpoint := fmt.Sprintf("%s%s", serverUrl, "/api/articles")
	httpClient := http.Client{}
	t.Run("Should list comments from an article", func(t *testing.T) {
		_, authorToken := integrationtests.MustRegisterUser(t, profileManagerRequests.RegisterPayload{})
		article := integrationtests.MustWriteArticle(t, articlePublisherRequests.WriteArticlePayload{}, authorToken)
		comment1 := integrationtests.MustWriteComment(t, articlePublisherRequests.WriteCommentPayload{}, article.Article.Slug, authorToken)
		comment2 := integrationtests.MustWriteComment(t, articlePublisherRequests.WriteCommentPayload{}, article.Article.Slug, authorToken)
		comment3 := integrationtests.MustWriteComment(t, articlePublisherRequests.WriteCommentPayload{}, article.Article.Slug, authorToken)
		comments := map[string]*articlePublisherResponses.CommentResponse{
			comment1.Comment.ID: comment1,
			comment2.Comment.ID: comment2,
			comment3.Comment.ID: comment3,
		}
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s/%s", listCommentsEndpoint, article.Article.Slug, commentsPath), nil)
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", authorToken))
		res, err := httpClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)
		resBytes, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		listCommentsResponse := new(articlePublisherResponses.CommentsResponse)
		err = json.Unmarshal(resBytes, listCommentsResponse)
		require.NoError(t, err)
		checkListCommentsResponse(t, listCommentsResponse, len(comments), comments)
	})
	t.Run("Should return http 404 if targeted article was not found", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s/%s", listCommentsEndpoint, uuid.NewString(), commentsPath), nil)
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res, err := httpClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusNotFound, res.StatusCode)
	})
}

func checkListCommentsResponse(t *testing.T, response *articlePublisherResponses.CommentsResponse, length int, comments map[string]*articlePublisherResponses.CommentResponse) {
	require.Len(t, response.Comment, length)
	for _, comment := range response.Comment {
		createdComment, ok := comments[comment.ID]
		require.True(t, ok)
		require.EqualValues(t, createdComment.Comment, comment)
	}
}
