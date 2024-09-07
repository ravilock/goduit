package articlepublisher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	integrationtests "github.com/ravilock/goduit/integrationTests"
	articlePublisherRequests "github.com/ravilock/goduit/internal/articlePublisher/requests"
	articlePublisherResponses "github.com/ravilock/goduit/internal/articlePublisher/responses"
	profileManagerRequests "github.com/ravilock/goduit/internal/profileManager/requests"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestWriteArticle(t *testing.T) {
	serverUrl := viper.GetString("server.url")
	writeArticleEndpoint := fmt.Sprintf("%s%s", serverUrl, "/api/articles")
	httpClient := http.Client{}
	t.Run("Should create an article", func(t *testing.T) {
		authorIdentity, authorToken := integrationtests.MustRegisterUser(t, profileManagerRequests.RegisterPayload{})
		createArticleRequest := generateWriteArticleBody()
		requestBody, err := json.Marshal(createArticleRequest)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodPost, writeArticleEndpoint, bytes.NewBuffer(requestBody))
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", authorToken))
		res, err := httpClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, res.StatusCode)
		resBytes, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		createArticleResponse := new(articlePublisherResponses.ArticleResponse)
		err = json.Unmarshal(resBytes, createArticleResponse)
		require.NoError(t, err)
		checkWriteArticleResponse(t, createArticleRequest, authorIdentity.Username, createArticleResponse)
	})
	t.Run("Should not allow to write article with slug that already exists", func(t *testing.T) {
		authorIdentity, authorToken := integrationtests.MustRegisterUser(t, profileManagerRequests.RegisterPayload{})
		createArticleRequest := generateWriteArticleBody()
		requestBody, err := json.Marshal(createArticleRequest)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodPost, writeArticleEndpoint, bytes.NewBuffer(requestBody))
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", authorToken))
		res, err := httpClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, res.StatusCode)
		resBytes, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		createArticleResponse := new(articlePublisherResponses.ArticleResponse)
		err = json.Unmarshal(resBytes, createArticleResponse)
		require.NoError(t, err)
		checkWriteArticleResponse(t, createArticleRequest, authorIdentity.Username, createArticleResponse)
		res, err = httpClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusConflict, res.StatusCode)
	})
}

func generateWriteArticleBody() *articlePublisherRequests.WriteArticleRequest {
	request := new(articlePublisherRequests.WriteArticleRequest)
	request.Article.Title = integrationtests.UniqueTitle()
	request.Article.Description = "Test Article Description"
	request.Article.Body = "Test Article Body"
	request.Article.TagList = []string{"test"}
	return request
}

func checkWriteArticleResponse(t *testing.T, request *articlePublisherRequests.WriteArticleRequest, author string, response *articlePublisherResponses.ArticleResponse) {
	t.Helper()
	require.Equal(t, request.Article.Title, response.Article.Title, "Wrong article title")
	require.Equal(t, request.Article.Description, response.Article.Description, "Wrong article description")
	require.Equal(t, request.Article.Body, response.Article.Body, "Wrong article body")
	require.Equal(t, integrationtests.MakeSlug(request.Article.Title), response.Article.Slug, "Wrong article body")
	require.Equal(t, request.Article.TagList, response.Article.TagList, "Wrong article body")
	require.Equal(t, author, response.Article.Author.Username, "Wrong article author username")
}
