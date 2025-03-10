package articlepublisher

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"testing"

	"github.com/labstack/echo/v4"
	integrationtests "github.com/ravilock/goduit/integrationTests"
	articlePublisherResponses "github.com/ravilock/goduit/internal/articlePublisher/responses"
	"github.com/ravilock/goduit/internal/mongo"
	profileManagerRequests "github.com/ravilock/goduit/internal/profileManager/requests"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestListArticles(t *testing.T) {
	client, err := mongo.ConnectDatabase(viper.GetString("db.url"))
	if err != nil {
		log.Fatalln("Error connecting to database", err)
	}
	integrationtests.ClearDatabase(client)
	serverUrl := viper.GetString("server.url")
	listArticlesEndpoint := fmt.Sprintf("%s%s", serverUrl, "/api/articles")
	httpClient := http.Client{}
	authorIdentity1, authorToken1 := integrationtests.MustRegisterUser(t, profileManagerRequests.RegisterPayload{})
	_, author1Tags := integrationtests.MustWriteArticles(t, 15, authorToken1, authorIdentity1.Username, authorIdentity1.Subject)
	authorIdentity2, authorToken2 := integrationtests.MustRegisterUser(t, profileManagerRequests.RegisterPayload{})
	_, author2Tags := integrationtests.MustWriteArticles(t, 15, authorToken2, authorIdentity2.Username, authorIdentity2.Subject)

	t.Run("Should list all articles", func(t *testing.T) {
		// Arrange
		limit := 30
		req, err := http.NewRequest(http.MethodGet, listArticlesEndpoint, nil)
		require.NoError(t, err)
		q := req.URL.Query()
		q.Add("limit", strconv.Itoa(limit))
		req.URL.RawQuery = q.Encode()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		// Act
		res, err := httpClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)
		resBytes, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		listArticleResponse := new(articlePublisherResponses.ArticlesResponse)
		err = json.Unmarshal(resBytes, listArticleResponse)
		require.NoError(t, err)
		checkListArticlesResponse(t, append(author1Tags, author2Tags...), []string{authorIdentity1.Username, authorIdentity2.Username}, limit, listArticleResponse)
	})

	t.Run("Should filter articles based on tag", func(t *testing.T) {
		// Arrange
		tag := authorIdentity1.Subject
		req, err := http.NewRequest(http.MethodGet, listArticlesEndpoint, nil)
		require.NoError(t, err)
		q := req.URL.Query()
		q.Add("tag", tag)
		req.URL.RawQuery = q.Encode()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		// Act
		res, err := httpClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)
		resBytes, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		listArticleResponse := new(articlePublisherResponses.ArticlesResponse)
		err = json.Unmarshal(resBytes, listArticleResponse)
		require.NoError(t, err)
		checkListArticlesResponse(t, author1Tags, []string{authorIdentity1.Username}, 15, listArticleResponse)
		for _, article := range listArticleResponse.Articles {
			require.NotContains(t, article.TagList, authorIdentity2.Subject)
			require.NotEqual(t, authorIdentity2.Username, article.Author.Username)
		}
	})

	t.Run("Should filter articles based on author", func(t *testing.T) {
		// Arrange
		author := authorIdentity2.Username
		req, err := http.NewRequest(http.MethodGet, listArticlesEndpoint, nil)
		require.NoError(t, err)
		q := req.URL.Query()
		q.Add("author", author)
		req.URL.RawQuery = q.Encode()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		// Act
		res, err := httpClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)
		resBytes, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		listArticleResponse := new(articlePublisherResponses.ArticlesResponse)
		err = json.Unmarshal(resBytes, listArticleResponse)
		require.NoError(t, err)
		checkListArticlesResponse(t, author2Tags, []string{authorIdentity2.Username}, 15, listArticleResponse)
		for _, article := range listArticleResponse.Articles {
			require.NotContains(t, article.TagList, authorIdentity1.Subject)
			require.NotEqual(t, authorIdentity1.Username, article.Author.Username)
		}
	})

	t.Run("Should limit the number of results", func(t *testing.T) {
		// Arrange
		limit := 3
		req, err := http.NewRequest(http.MethodGet, listArticlesEndpoint, nil)
		require.NoError(t, err)
		q := req.URL.Query()
		q.Add("limit", strconv.Itoa(limit))
		req.URL.RawQuery = q.Encode()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		// Act
		res, err := httpClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)
		resBytes, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		listArticleResponse := new(articlePublisherResponses.ArticlesResponse)
		err = json.Unmarshal(resBytes, listArticleResponse)
		require.NoError(t, err)
		checkListArticlesResponse(t, append(author1Tags, author2Tags...), []string{authorIdentity1.Username, authorIdentity2.Username}, limit, listArticleResponse)
	})

	t.Run("Should properly offset results", func(t *testing.T) {
		// Arrange
		offset := 5
		limit := 15
		req, err := http.NewRequest(http.MethodGet, listArticlesEndpoint, nil)
		require.NoError(t, err)
		q := req.URL.Query()
		q.Add("limit", strconv.Itoa(limit))
		req.URL.RawQuery = q.Encode()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		// Act
		res, err := httpClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)
		resBytes, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		listArticlesResponse1 := new(articlePublisherResponses.ArticlesResponse)
		err = json.Unmarshal(resBytes, listArticlesResponse1)
		require.NoError(t, err)
		req, err = http.NewRequest(http.MethodGet, listArticlesEndpoint, nil)
		require.NoError(t, err)
		q = req.URL.Query()
		q.Add("limit", strconv.Itoa(limit))
		q.Add("offset", strconv.Itoa(offset))
		req.URL.RawQuery = q.Encode()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		// Act
		res, err = httpClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)
		resBytes, err = io.ReadAll(res.Body)
		require.NoError(t, err)
		listArticlesResponse2 := new(articlePublisherResponses.ArticlesResponse)
		err = json.Unmarshal(resBytes, listArticlesResponse2)
		require.NoError(t, err)

		// Check that the responses intersect exactly on 10 elements
		articles1 := listArticlesResponse1.Articles[5:16]
		articles2 := listArticlesResponse2.Articles[:11]
		for i := 0; i < limit-offset; i++ {
			article1 := &articles1[i]
			article2 := &articles2[i]
			require.True(t, checkArticlesAreTheSame(article1, article2))
		}
	})
}

func checkListArticlesResponse(t *testing.T, tags []string, authorsUsername []string, limit int, response *articlePublisherResponses.ArticlesResponse) {
	t.Helper()
	for _, article := range response.Articles {
		require.Equal(t, len(response.Articles), limit)
		for _, tag := range tags {
			require.Contains(t, article.TagList, tag)
		}
		require.Contains(t, authorsUsername, article.Author.Username)
	}
}

func checkArticlesAreTheSame(articleA, articleB *articlePublisherResponses.MultiArticle) bool {
	return articleA.Slug == articleB.Slug
}
