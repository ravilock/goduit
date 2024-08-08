package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/api/validators"
	"github.com/ravilock/goduit/internal/app"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUnpublishArticle(t *testing.T) {
	err := validators.InitValidator()
	require.NoError(t, err)
	articleUnpublisherMock := newMockArticleUnpublisher(t)
	handler := unpublishArticleHandler{articleUnpublisherMock}
	e := echo.New()

	t.Run("Should delete an article", func(t *testing.T) {
		// Arrange
		articleAuthorID := primitive.NewObjectID()
		expectedAuthor := assembleArticleAuthor(articleAuthorID.Hex())
		expectedArticle := assembleArticleModel(articleAuthorID)
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/article/%s", *expectedArticle.Slug), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", expectedAuthor.ID.Hex())
		req.Header.Set("Goduit-Client-Username", *expectedAuthor.Username)
		req.Header.Set("Goduit-Client-Email", *expectedAuthor.Email)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("slug")
		c.SetParamValues(*expectedArticle.Slug)
		ctx := c.Request().Context()
		articleUnpublisherMock.EXPECT().GetArticleBySlug(ctx, *expectedArticle.Slug).Return(expectedArticle, nil).Once()
		articleUnpublisherMock.EXPECT().UnpublishArticle(ctx, *expectedArticle.Slug).Return(nil).Once()

		// Act
		err := handler.UnpublishArticle(c)

		// Assert
		require.NoError(t, err)
		require.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("Should return HTTP 404 if no article is found", func(t *testing.T) {
		// Arrange
		articleAuthorID := primitive.NewObjectID()
		expectedAuthor := assembleArticleAuthor(articleAuthorID.Hex())
		expectedArticle := assembleArticleModel(articleAuthorID)
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/article/%s", *expectedArticle.Slug), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", expectedAuthor.ID.Hex())
		req.Header.Set("Goduit-Client-Username", *expectedAuthor.Username)
		req.Header.Set("Goduit-Client-Email", *expectedAuthor.Email)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("slug")
		c.SetParamValues(*expectedArticle.Slug)
		articleUnpublisherMock.EXPECT().GetArticleBySlug(c.Request().Context(), *expectedArticle.Slug).Return(nil, app.ArticleNotFoundError(*expectedArticle.Slug, nil)).Once()

		// Act
		err := handler.UnpublishArticle(c)

		// Assert
		require.ErrorContains(t, err, api.ArticleNotFound(*expectedArticle.Slug).Error())
	})

	t.Run("Should only delete aritcles authored by the currently authenticated user", func(t *testing.T) {
		// Arrange
		articleAuthorID := primitive.NewObjectID()
		expectedArticle := assembleArticleModel(articleAuthorID)
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/article/%s", *expectedArticle.Slug), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", uuid.NewString())
		req.Header.Set("Goduit-Client-Username", "not-the-author")
		req.Header.Set("Goduit-Client-Email", "not.the.author.email@test.test")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("slug")
		c.SetParamValues(*expectedArticle.Slug)
		ctx := c.Request().Context()
		articleUnpublisherMock.EXPECT().GetArticleBySlug(ctx, *expectedArticle.Slug).Return(expectedArticle, nil).Once()

		// Act
		err := handler.UnpublishArticle(c)

		// Assert
		require.ErrorContains(t, err, api.Forbidden.Error())
	})
}
