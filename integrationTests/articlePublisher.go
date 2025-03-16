package integrationtests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	articlePublisherModels "github.com/ravilock/goduit/internal/articlePublisher/models"
	articlePublisherRepositories "github.com/ravilock/goduit/internal/articlePublisher/repositories"
	articlePublisherRequests "github.com/ravilock/goduit/internal/articlePublisher/requests"
	articlePublisherResponses "github.com/ravilock/goduit/internal/articlePublisher/responses"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
)

func MakeSlug(title string) string {
	loweredTitle := strings.ToLower(title)
	return strings.ReplaceAll(loweredTitle, " ", "-")
}

func MustWriteArticle(t *testing.T, writeArticlePayload articlePublisherRequests.WriteArticlePayload, cookie *http.Cookie) *articlePublisherResponses.ArticleResponse {
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
	req.AddCookie(cookie)
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

func GenerateArticleModel(authorID string) *articlePublisherModels.Article {
	title := UniqueTitle()
	slug := MakeSlug(title)
	description := "Article Description"
	body := "Article Body"
	tagList := []string{"categories", "housing", "technology"}
	now := time.Now()
	return &articlePublisherModels.Article{
		Author:         &authorID,
		Slug:           &slug,
		Title:          &title,
		Description:    &description,
		Body:           &body,
		TagList:        tagList,
		CreatedAt:      &now,
		UpdatedAt:      &now,
		FavoritesCount: new(int64),
	}
}

func MustWriteArticleRegister(t *testing.T, client *mongo.Client, articleModel *articlePublisherModels.Article) {
	articleRepository := articlePublisherRepositories.NewArticleRepository(client)
	err := articleRepository.WriteArticle(context.Background(), articleModel)
	require.NoError(t, err)
}

func MustWriteArticles(t *testing.T, amount int, authorCookie *http.Cookie, authorUsername, authorID string) ([]*articlePublisherResponses.ArticleResponse, []string) {
	articles := make([]*articlePublisherResponses.ArticleResponse, 0, amount)
	tags := make([]string, 0, amount)
	for i := 0; i < amount; i++ {
		article := MustWriteArticle(t, articlePublisherRequests.WriteArticlePayload{TagList: []string{authorID}}, authorCookie)
		articles = append(articles, article)
	}
	return articles, slices.Compact(tags)
}

func MustWriteComment(t *testing.T, writeCommentPayload articlePublisherRequests.WriteCommentPayload, articleSlug string, authorCookie *http.Cookie) *articlePublisherResponses.CommentResponse {
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
	req.AddCookie(authorCookie)
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
