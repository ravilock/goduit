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

func TestGetArticle(t *testing.T) {
	serverUrl := viper.GetString("server.url")
	publicArticleEndpoint := fmt.Sprintf("%s%s", serverUrl, "/api/articles")
	httpClient := http.Client{}
	t.Run("Should get an article", func(t *testing.T) {
		authorIdentity, authorToken := integrationtests.MustRegisterUser(t, profileManagerRequests.RegisterPayload{})
		article := integrationtests.MustWriteArticle(t, articlePublisherRequests.WriteArticlePayload{}, authorToken)
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", publicArticleEndpoint, article.Article.Slug), nil)
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", authorToken))
		res, err := httpClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)
		resBytes, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		getArticleResponse := new(articlePublisherResponses.ArticleResponse)
		err = json.Unmarshal(resBytes, getArticleResponse)
		require.NoError(t, err)
		checkGetArticleResponse(t, article.Article.Title, article.Article.Slug, authorIdentity.Username, article.Article.TagList, getArticleResponse)
	})
	t.Run("Should return HTTP 404 if no article is found", func(t *testing.T) {
		_, authorToken := integrationtests.MustRegisterUser(t, profileManagerRequests.RegisterPayload{})
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", publicArticleEndpoint, uuid.NewString()), nil)
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", authorToken))
		res, err := httpClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusNotFound, res.StatusCode)
	})
}

func checkGetArticleResponse(t *testing.T, title, slug, authorUsername string, tagList []string, response *articlePublisherResponses.ArticleResponse) {
	t.Helper()
	require.Equal(t, authorUsername, response.Article.Author.Username, "Article's author username is wrong")
	require.Equal(t, title, response.Article.Title, "Wrong article title")
	require.Equal(t, slug, response.Article.Slug, "Wrong article slug")
	require.Equal(t, tagList, response.Article.TagList, "Wrong article tag list")
}
