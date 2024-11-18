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

func TestUpdateArticle(t *testing.T) {
	serverUrl := viper.GetString("server.url")
	updateArticleEndpoint := fmt.Sprintf("%s%s", serverUrl, "/api/articles")
	httpClient := http.Client{}

	t.Run("Should update an article", func(t *testing.T) {
		// Arrange
		authorIdentity, authorToken := integrationtests.MustRegisterUser(t, profileManagerRequests.RegisterPayload{})
		article := integrationtests.MustWriteArticle(t, articlePublisherRequests.WriteArticlePayload{}, authorToken)
		updateArticleRequest := generateUpdateArticleBody()
		requestBody, err := json.Marshal(updateArticleRequest)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%s", updateArticleEndpoint, article.Article.Slug), bytes.NewBuffer(requestBody))
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", authorToken))

		// Act
		res, err := httpClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)
		resBytes, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		updateArticleResponse := new(articlePublisherResponses.ArticleResponse)
		err = json.Unmarshal(resBytes, updateArticleResponse)
		require.NoError(t, err)
		checkUpdateArticleResponse(t, updateArticleRequest, authorIdentity.Username, updateArticleResponse, article.Article.TagList)
	})

	t.Run("Should be able to update article keeping it's title", func(t *testing.T) {
		// Arrange
		authorIdentity, authorToken := integrationtests.MustRegisterUser(t, profileManagerRequests.RegisterPayload{})
		article := integrationtests.MustWriteArticle(t, articlePublisherRequests.WriteArticlePayload{}, authorToken)
		updateArticleRequest := generateUpdateArticleBody()
		updateArticleRequest.Article.Title = article.Article.Title
		requestBody, err := json.Marshal(updateArticleRequest)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%s", updateArticleEndpoint, article.Article.Slug), bytes.NewBuffer(requestBody))
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", authorToken))

		// Act
		res, err := httpClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)
		resBytes, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		updateArticleResponse := new(articlePublisherResponses.ArticleResponse)
		err = json.Unmarshal(resBytes, updateArticleResponse)
		require.NoError(t, err)
		checkUpdateArticleResponse(t, updateArticleRequest, authorIdentity.Username, updateArticleResponse, article.Article.TagList)
	})

	t.Run("Should not allow to update article's slug to another that already exists", func(t *testing.T) {
		// Arrange
		_, authorToken := integrationtests.MustRegisterUser(t, profileManagerRequests.RegisterPayload{})
		article := integrationtests.MustWriteArticle(t, articlePublisherRequests.WriteArticlePayload{}, authorToken)
		conflictedArticle := integrationtests.MustWriteArticle(t, articlePublisherRequests.WriteArticlePayload{}, authorToken)
		updateArticleRequest := generateUpdateArticleBody()
		updateArticleRequest.Article.Title = conflictedArticle.Article.Title
		requestBody, err := json.Marshal(updateArticleRequest)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%s", updateArticleEndpoint, article.Article.Slug), bytes.NewBuffer(requestBody))
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", authorToken))

		// Act
		res, err := httpClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusConflict, res.StatusCode)
	})

	t.Run("Should return HTTP 404 if targeted article does not exists", func(t *testing.T) {
		// Arrange
		_, authorToken := integrationtests.MustRegisterUser(t, profileManagerRequests.RegisterPayload{})
		updateArticleRequest := generateUpdateArticleBody()
		requestBody, err := json.Marshal(updateArticleRequest)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%s", updateArticleEndpoint, uuid.NewString()), bytes.NewBuffer(requestBody))
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", authorToken))

		// Act
		res, err := httpClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusNotFound, res.StatusCode)
	})

	t.Run("Should not allow user to update other author's articles", func(t *testing.T) {
		// Arrange
		_, authorToken := integrationtests.MustRegisterUser(t, profileManagerRequests.RegisterPayload{})
		article := integrationtests.MustWriteArticle(t, articlePublisherRequests.WriteArticlePayload{}, authorToken)
		_, nonAuthorToken := integrationtests.MustRegisterUser(t, profileManagerRequests.RegisterPayload{})
		updateArticleRequest := generateUpdateArticleBody()
		requestBody, err := json.Marshal(updateArticleRequest)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%s", updateArticleEndpoint, article.Article.Slug), bytes.NewBuffer(requestBody))
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", nonAuthorToken))
		res, err := httpClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusForbidden, res.StatusCode)
	})
}

func generateUpdateArticleBody() *articlePublisherRequests.UpdateArticleRequest {
	article := new(articlePublisherRequests.UpdateArticleRequest)
	article.Article.Title = integrationtests.UniqueTitle()
	article.Article.Description = "Test Description"
	article.Article.Body = "Test Body"
	return article
}

func checkUpdateArticleResponse(t *testing.T, request *articlePublisherRequests.UpdateArticleRequest, authorUsername string, response *articlePublisherResponses.ArticleResponse, tagList []string) {
	t.Helper()
	require.Equal(t, request.Article.Title, response.Article.Title, "Wrong article title")
	require.Equal(t, request.Article.Description, response.Article.Description, "Wrong article description")
	require.Equal(t, request.Article.Body, response.Article.Body, "Wrong article body")
	require.Equal(t, integrationtests.MakeSlug(request.Article.Title), response.Article.Slug, "Wrong article body")
	require.Equal(t, tagList, response.Article.TagList, "Wrong article body")
	require.Equal(t, authorUsername, response.Article.Author.Username, "Wrong article author username")
}
