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

func TestWriteComment(t *testing.T) {
	err := validators.InitValidator()
	require.NoError(t, err)
	commentWriterMock := newMockCommentWriter(t)
	articleGetterMock := newMockArticleGetter(t)
	profileGetterMock := newMockProfileGetter(t)
	handler := &writeCommentHandler{commentWriterMock, articleGetterMock, profileGetterMock}

	e := echo.New()

	t.Run("Should create a commentary", func(t *testing.T) {
		// Arrange
		articleAuthorID := primitive.NewObjectID()
		expectedAuthor := assembleArticleAuthor(articleAuthorID.Hex())
		expectedArticle := assembleArticleModel(articleAuthorID)
		createCommentRequest := generateWriteCommentBody()
		expectedCommentModel := createCommentRequest.Model(articleAuthorID.Hex())
		articleID := expectedArticle.ID.Hex()
		expectedCommentModel.Article = &articleID
		requestBody, err := json.Marshal(createCommentRequest)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/articles/%s/comments", *expectedArticle.Slug), bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", expectedAuthor.ID.Hex())
		req.Header.Set("Goduit-Client-Username", *expectedAuthor.Username)
		req.Header.Set("Goduit-Client-Email", *expectedAuthor.Email)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("slug")
		c.SetParamValues(*expectedArticle.Slug)
		ctx := c.Request().Context()
		articleGetterMock.EXPECT().GetArticleBySlug(ctx, *expectedArticle.Slug).Return(expectedArticle, nil).Once()
		profileGetterMock.EXPECT().GetProfileByID(ctx, expectedAuthor.ID.Hex()).Return(expectedAuthor, nil).Once()
		commentWriterMock.EXPECT().WriteComment(ctx, expectedCommentModel).RunAndReturn(func(ctx context.Context, comment *models.Comment) error {
			commentID := primitive.NewObjectID()
			comment.ID = &commentID
			return nil
		}).Once()

		// Act
		err = handler.WriteComment(c)

		// Assert
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, rec.Code)
		createCommentResponse := new(articlePublisherResponses.CommentResponse)
		err = json.Unmarshal(rec.Body.Bytes(), createCommentResponse)
		require.NoError(t, err)
		checkWriteCommentResponse(t, createCommentRequest, *expectedAuthor.Username, createCommentResponse)
	})

	t.Run("Should return HTTP 404 if no article is found", func(t *testing.T) {
		// Arrange
		articleAuthorID := primitive.NewObjectID()
		expectedAuthor := assembleArticleAuthor(articleAuthorID.Hex())
		expectedArticle := assembleArticleModel(articleAuthorID)
		createCommentRequest := generateWriteCommentBody()
		requestBody, err := json.Marshal(createCommentRequest)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/articles/%s/comments", *expectedArticle.Slug), bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", expectedAuthor.ID.Hex())
		req.Header.Set("Goduit-Client-Username", *expectedAuthor.Username)
		req.Header.Set("Goduit-Client-Email", *expectedAuthor.Email)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("slug")
		c.SetParamValues(*expectedArticle.Slug)
		articleGetterMock.EXPECT().GetArticleBySlug(c.Request().Context(), *expectedArticle.Slug).Return(nil, app.ArticleNotFoundError(*expectedArticle.Slug, nil)).Once()

		// Act
		err = handler.WriteComment(c)

		// Assert
		require.ErrorContains(t, err, api.ArticleNotFound(*expectedArticle.Slug).Error())
	})
}

func generateWriteCommentBody() *articlePublisherRequests.WriteCommentRequest {
	request := new(articlePublisherRequests.WriteCommentRequest)
	request.Comment.Body = "Test Comment Body"
	return request
}

func checkWriteCommentResponse(t *testing.T, request *articlePublisherRequests.WriteCommentRequest, author string, response *articlePublisherResponses.CommentResponse) {
	t.Helper()
	require.NotZero(t, response.Comment.ID)
	require.Equal(t, request.Comment.Body, response.Comment.Body, "Wrong comment body")
	require.Equal(t, author, response.Comment.Author.Username, "Wrong comment author username")
}
