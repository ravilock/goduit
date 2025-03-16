package articlepublisher

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	integrationtests "github.com/ravilock/goduit/integrationTests"
	articlePublisherRequests "github.com/ravilock/goduit/internal/articlePublisher/requests"
	profileManagerRequests "github.com/ravilock/goduit/internal/profileManager/requests"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestUnpublishArticle(t *testing.T) {
	serverUrl := viper.GetString("server.url")
	unpublishArticleEndpoint := fmt.Sprintf("%s%s", serverUrl, "/api/articles")
	httpClient := http.Client{}

	t.Run("Should unpublish an article", func(t *testing.T) {
		// Arrange
		_, authorCookie := integrationtests.MustRegisterUser(t, profileManagerRequests.RegisterPayload{})
		article := integrationtests.MustWriteArticle(t, articlePublisherRequests.WriteArticlePayload{}, authorCookie)
		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s", unpublishArticleEndpoint, article.Article.Slug), nil)
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.AddCookie(authorCookie)

		// Act
		res, err := httpClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusNoContent, res.StatusCode)
	})

	t.Run("Should return HTTP 404 if targeted article does not exists", func(t *testing.T) {
		// Arrange
		_, authorCookie := integrationtests.MustRegisterUser(t, profileManagerRequests.RegisterPayload{})
		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s", unpublishArticleEndpoint, uuid.NewString()), nil)
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.AddCookie(authorCookie)

		// Act
		res, err := httpClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusNotFound, res.StatusCode)
	})

	t.Run("Should not allow user to delete other author's articles", func(t *testing.T) {
		// Arrange
		_, authorCookie := integrationtests.MustRegisterUser(t, profileManagerRequests.RegisterPayload{})
		article := integrationtests.MustWriteArticle(t, articlePublisherRequests.WriteArticlePayload{}, authorCookie)
		_, nonAuthorCookie := integrationtests.MustRegisterUser(t, profileManagerRequests.RegisterPayload{})
		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s", unpublishArticleEndpoint, article.Article.Slug), nil)
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.AddCookie(nonAuthorCookie)

		// Act
		res, err := httpClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusForbidden, res.StatusCode)
	})
}
