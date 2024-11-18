package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/api/validators"
	"github.com/ravilock/goduit/internal/app"
	"github.com/ravilock/goduit/internal/articlePublisher/models"
	articlePublisherRequests "github.com/ravilock/goduit/internal/articlePublisher/requests"
	articlePublisherResponses "github.com/ravilock/goduit/internal/articlePublisher/responses"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUpdateArticle(t *testing.T) {
	err := validators.InitValidator()
	require.NoError(t, err)
	articleUpdaterMock := newMockArticleUpdater(t)
	profileGetterMock := newMockProfileGetter(t)
	handler := &updateArticleHandler{articleUpdaterMock, profileGetterMock}

	e := echo.New()

	t.Run("Should update an article", func(t *testing.T) {
		// Arrange
		expectedAuthor := assembleRandomUser()
		expectedArticle := assembleArticleModel(*expectedAuthor.ID)
		updateArticleRequest := generateUpdateArticleBody()
		requestBody, err := json.Marshal(updateArticleRequest)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/article/%s", *expectedArticle.Slug), bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", expectedAuthor.ID.Hex())
		req.Header.Set("Goduit-Client-Username", *expectedAuthor.Username)
		req.Header.Set("Goduit-Client-Email", *expectedAuthor.Email)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("slug")
		c.SetParamValues(*expectedArticle.Slug)
		ctx := c.Request().Context()
		articleUpdaterMock.EXPECT().GetArticleBySlug(ctx, *expectedArticle.Slug).Return(expectedArticle, nil).Once()
		articleUpdaterMock.EXPECT().UpdateArticle(ctx, *expectedArticle.Slug, updateArticleRequest.Model()).RunAndReturn(func(ctx context.Context, slug string, article *models.Article) error {
			favoritesCount := int64(30)
			article.FavoritesCount = &favoritesCount
			article.TagList = expectedArticle.TagList
			return nil
		}).Once()
		profileGetterMock.EXPECT().GetProfileByID(ctx, expectedAuthor.ID.Hex()).Return(expectedAuthor, nil).Once()

		// Act
		err = handler.UpdateArticle(c)

		// Assert
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rec.Code)
		updateArticleResponse := new(articlePublisherResponses.ArticleResponse)
		err = json.Unmarshal(rec.Body.Bytes(), updateArticleResponse)
		require.NoError(t, err)
		checkUpdateArticleResponse(t, updateArticleRequest, *expectedAuthor.Username, updateArticleResponse, expectedArticle.TagList)
	})

	t.Run("Should return HTTP 404 if no article is found", func(t *testing.T) {
		// Arrange
		expectedAuthor := assembleRandomUser()
		updateArticleRequest := generateUpdateArticleBody()
		articleSlug := "test-article"
		requestBody, err := json.Marshal(updateArticleRequest)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/article/%s", articleSlug), bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", expectedAuthor.ID.Hex())
		req.Header.Set("Goduit-Client-Username", *expectedAuthor.Username)
		req.Header.Set("Goduit-Client-Email", *expectedAuthor.Email)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("slug")
		c.SetParamValues(articleSlug)
		ctx := c.Request().Context()
		articleUpdaterMock.EXPECT().GetArticleBySlug(ctx, articleSlug).Return(nil, app.ArticleNotFoundError(articleSlug, nil)).Once()

		// Act
		err = handler.UpdateArticle(c)

		// Assert
		require.ErrorContains(t, err, api.ArticleNotFound(articleSlug).Error())
	})

	t.Run("Should only update articles authored by the currently authenticaed user", func(t *testing.T) {
		// Arrange
		expectedAuthor := assembleRandomUser()
		expectedArticle := assembleArticleModel(*expectedAuthor.ID)
		updateArticleRequest := generateUpdateArticleBody()
		requestBody, err := json.Marshal(updateArticleRequest)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/article/%s", *expectedArticle.Slug), bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", primitive.NewObjectID().Hex())
		req.Header.Set("Goduit-Client-Username", "not-the-author")
		req.Header.Set("Goduit-Subject", "not.the.author.email@test.test")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("slug")
		c.SetParamValues(*expectedArticle.Slug)
		ctx := c.Request().Context()
		articleUpdaterMock.EXPECT().GetArticleBySlug(ctx, *expectedArticle.Slug).Return(expectedArticle, nil).Once()

		// Act
		require.NoError(t, err)

		// Assert
		err = handler.UpdateArticle(c)
		require.ErrorContains(t, err, api.Forbidden.Error())
	})
}

func generateUpdateArticleBody() *articlePublisherRequests.UpdateArticleRequest {
	request := new(articlePublisherRequests.UpdateArticleRequest)
	request.Article.Title = "New Article Name"
	request.Article.Description = "New Article Description"
	request.Article.Body = "New Article Body"
	return request
}

func checkUpdateArticleResponse(t *testing.T, request *articlePublisherRequests.UpdateArticleRequest, author string, response *articlePublisherResponses.ArticleResponse, tagList []string) {
	t.Helper()
	require.Equal(t, request.Article.Title, response.Article.Title, "Wrong article title")
	require.Equal(t, request.Article.Description, response.Article.Description, "Wrong article description")
	require.Equal(t, request.Article.Body, response.Article.Body, "Wrong article body")
	require.Equal(t, makeSlug(request.Article.Title), response.Article.Slug, "Wrong article body")
	require.Equal(t, tagList, response.Article.TagList, "Wrong article body")
	require.Equal(t, author, response.Article.Author.Username, "Wrong article author username")
}
