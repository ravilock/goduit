package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/api/validators"
	"github.com/ravilock/goduit/internal/app"
	articlePublisherRequests "github.com/ravilock/goduit/internal/articlePublisher/requests"
	articlePublisherResponses "github.com/ravilock/goduit/internal/articlePublisher/responses"
	"github.com/stretchr/testify/require"
)

func TestWriteArticle(t *testing.T) {
	validators.InitValidator()
	articleWriterMock := newMockArticleWriter(t)
	profileGetterMock := newMockProfileGetter(t)
	handler := &writeArticleHandler{articleWriterMock, profileGetterMock}

	e := echo.New()

	t.Run("Should create an article", func(t *testing.T) {
		// Arrange
		expectedAuthor := assembleRandomUser()
		createArticleRequest := generateWriteArticleBody()
		requestBody, err := json.Marshal(createArticleRequest)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/api/articles", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", expectedAuthor.ID.Hex())
		req.Header.Set("Goduit-Client-Username", *expectedAuthor.Username)
		req.Header.Set("Goduit-Client-Email", *expectedAuthor.Email)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		ctx := c.Request().Context()
		articleWriterMock.EXPECT().WriteArticle(ctx, createArticleRequest.Model(expectedAuthor.ID.Hex())).Return(nil).Once()
		profileGetterMock.EXPECT().GetProfileByID(ctx, expectedAuthor.ID.Hex()).Return(expectedAuthor, nil).Once()

		// Act
		err = handler.WriteArticle(c)

		// Assert
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, rec.Code)
		createArticleResponse := new(articlePublisherResponses.ArticleResponse)
		err = json.Unmarshal(rec.Body.Bytes(), createArticleResponse)
		require.NoError(t, err)
		checkWriteArticleResponse(t, createArticleRequest, *expectedAuthor.Username, createArticleResponse)
	})

	t.Run("Should return conflict error if title/slug is already used", func(t *testing.T) {
		// Arrange
		expectedAuthor := assembleRandomUser()
		createArticleRequest := generateWriteArticleBody()
		requestBody, err := json.Marshal(createArticleRequest)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/api/articles", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", expectedAuthor.ID.Hex())
		req.Header.Set("Goduit-Client-Username", *expectedAuthor.Username)
		req.Header.Set("Goduit-Client-Email", *expectedAuthor.Email)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		ctx := c.Request().Context()
		articleWriterMock.EXPECT().WriteArticle(ctx, createArticleRequest.Model(expectedAuthor.ID.Hex())).Return(app.ConflictError("articles")).Once()

		// Act
		err = handler.WriteArticle(c)

		// Assert
		require.ErrorIs(t, err, api.ConfictError)
	})
}

func generateWriteArticleBody() *articlePublisherRequests.WriteArticleRequest {
	request := new(articlePublisherRequests.WriteArticleRequest)
	request.Article.Title = "Test Article Name"
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
	require.Equal(t, makeSlug(request.Article.Title), response.Article.Slug, "Wrong article body")
	require.Equal(t, request.Article.TagList, response.Article.TagList, "Wrong article body")
	require.Equal(t, author, response.Article.Author.Username, "Wrong article author username")
}
