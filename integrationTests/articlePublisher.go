package integrationtests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	articlePublisherRequests "github.com/ravilock/goduit/internal/articlePublisher/requests"
	articlePublisherResponses "github.com/ravilock/goduit/internal/articlePublisher/responses"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func MakeSlug(title string) string {
	loweredTitle := strings.ToLower(title)
	return strings.ReplaceAll(loweredTitle, " ", "-")
}

func MustWriteArticle(t *testing.T, writeArticlePayload articlePublisherRequests.WriteArticlePayload, token string) *articlePublisherResponses.ArticleResponse {
	httpClient := http.Client{}
	serverUrl := viper.GetString("server.url")
	writeArticleEndpoint := fmt.Sprintf("%s%s", serverUrl, "/api/articles")
	if writeArticlePayload.Body == "" {
		writeArticlePayload.Body = "Article Body"
	}
	if writeArticlePayload.Title == "" {
		writeArticlePayload.Title = UniqueTitle()
	}
	if len(writeArticlePayload.TagList) == 0 {
		writeArticlePayload.TagList = []string{"categories", "housing", "technology"}
	} else {
		writeArticlePayload.TagList = append(writeArticlePayload.TagList, "categories", "housing", "technology")
	}
	if writeArticlePayload.Description == "" {
		writeArticlePayload.Description = "Article Description"
	}
	requestBody, err := json.Marshal(&articlePublisherRequests.WriteArticleRequest{
		Article: writeArticlePayload,
	})
	require.NoError(t, err)
	req, err := http.NewRequest(http.MethodPost, writeArticleEndpoint, bytes.NewBuffer(requestBody))
	require.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))
	res, err := httpClient.Do(req)
	require.NoError(t, err)
	defer res.Body.Close()
	writeArticleResponse := new(articlePublisherResponses.ArticleResponse)
	resBytes, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.StatusCode)
	err = json.Unmarshal(resBytes, writeArticleResponse)
	require.NoError(t, err)
	return writeArticleResponse
}

func MustWriteArticles(t *testing.T, amount int, authorToken, authorUsername, authorID string) ([]*articlePublisherResponses.ArticleResponse, []string) {
	articles := make([]*articlePublisherResponses.ArticleResponse, 0, amount)
	tags := make([]string, 0, amount)
	for i := 0; i < amount; i++ {
		article := MustWriteArticle(t, articlePublisherRequests.WriteArticlePayload{TagList: []string{authorID}}, authorToken)
		articles = append(articles, article)
	}
	return articles, slices.Compact(tags)
}

func MustWriteComment(t *testing.T, writeCommentPayload articlePublisherRequests.WriteCommentPayload, articleSlug, token string) *articlePublisherResponses.CommentResponse {
	httpClient := http.Client{}
	serverUrl := viper.GetString("server.url")
	writeArticleEndpoint := fmt.Sprintf("%s%s/%s%s", serverUrl, "/api/articles", articleSlug, "/comments")
	if writeCommentPayload.Body == "" {
		writeCommentPayload.Body = UniqueTitle()
	}
	requestBody, err := json.Marshal(&articlePublisherRequests.WriteCommentRequest{
		Comment: writeCommentPayload,
	})
	require.NoError(t, err)
	req, err := http.NewRequest(http.MethodPost, writeArticleEndpoint, bytes.NewBuffer(requestBody))
	require.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))
	res, err := httpClient.Do(req)
	require.NoError(t, err)
	defer res.Body.Close()
	writeCommentResponse := new(articlePublisherResponses.CommentResponse)
	resBytes, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.StatusCode)
	err = json.Unmarshal(resBytes, writeCommentResponse)
	require.NoError(t, err)
	return writeCommentResponse
}
