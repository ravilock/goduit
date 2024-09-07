package articlepublisher

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	integrationtests "github.com/ravilock/goduit/integrationTests"
	articlePublisherRequests "github.com/ravilock/goduit/internal/articlePublisher/requests"
	profileManagerRequests "github.com/ravilock/goduit/internal/profileManager/requests"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDeleteComment(t *testing.T) {
	serverUrl := viper.GetString("server.url")
	deleteCommentEndpoint := fmt.Sprintf("%s%s", serverUrl, "/api/articles")
	httpClient := http.Client{}
	t.Run("Should delete a comment", func(t *testing.T) {
		_, authorToken := integrationtests.MustRegisterUser(t, profileManagerRequests.RegisterPayload{})
		article := integrationtests.MustWriteArticle(t, articlePublisherRequests.WriteArticlePayload{}, authorToken)
		comment := integrationtests.MustWriteComment(t, articlePublisherRequests.WriteCommentPayload{}, article.Article.Slug, authorToken)
		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s/%s/%s", deleteCommentEndpoint, article.Article.Slug, commentsPath, comment.Comment.ID), nil)
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", authorToken))
		res, err := httpClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusNoContent, res.StatusCode)
	})
	t.Run("Should not allow user to delete other author's comments", func(t *testing.T) {
		_, authorToken := integrationtests.MustRegisterUser(t, profileManagerRequests.RegisterPayload{})
		article := integrationtests.MustWriteArticle(t, articlePublisherRequests.WriteArticlePayload{}, authorToken)
		comment := integrationtests.MustWriteComment(t, articlePublisherRequests.WriteCommentPayload{}, article.Article.Slug, authorToken)
		_, nonAuthorToken := integrationtests.MustRegisterUser(t, profileManagerRequests.RegisterPayload{})
		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s/%s/%s", deleteCommentEndpoint, article.Article.Slug, commentsPath, comment.Comment.ID), nil)
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", nonAuthorToken))
		res, err := httpClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusForbidden, res.StatusCode)
	})
	t.Run("Should return HTTP 404 if targeted article does not exists", func(t *testing.T) {
		_, authorToken := integrationtests.MustRegisterUser(t, profileManagerRequests.RegisterPayload{})
		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s/%s/%s", deleteCommentEndpoint, uuid.NewString(), commentsPath, uuid.NewString()), nil)
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", authorToken))
		res, err := httpClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusNotFound, res.StatusCode)
	})
	t.Run("Should return HTTP 404 if targeted comment does not exists", func(t *testing.T) {
		_, authorToken := integrationtests.MustRegisterUser(t, profileManagerRequests.RegisterPayload{})
		article := integrationtests.MustWriteArticle(t, articlePublisherRequests.WriteArticlePayload{}, authorToken)
		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s/%s/%s", deleteCommentEndpoint, article.Article.Slug, commentsPath, primitive.NewObjectID().Hex()), nil)
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", authorToken))
		res, err := httpClient.Do(req)
		require.NoError(t, err)
		resBytes, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		fmt.Println(string(resBytes))
		require.Equal(t, http.StatusNotFound, res.StatusCode)
	})
}
