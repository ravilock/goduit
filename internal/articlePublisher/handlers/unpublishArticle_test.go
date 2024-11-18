package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/api/validators"
	"github.com/ravilock/goduit/internal/app"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUnpublishArticle(t *testing.T) {
	validators.InitValidator()
	articleUnpublisherMock := newMockArticleUnpublisher(t)
	handler := unpublishArticleHandler{articleUnpublisherMock}
	e := echo.New()

	t.Run("Should delete an article", func(t *testing.T) {
		// Arrange
		expectedAuthor := assembleRandomUser()
		expectedArticle := assembleArticleModel(*expectedAuthor.ID)
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
		expectedAuthor := assembleRandomUser()
		slug := "inexistent-article"
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/article/%s", slug), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", expectedAuthor.ID.Hex())
		req.Header.Set("Goduit-Client-Username", *expectedAuthor.Username)
		req.Header.Set("Goduit-Client-Email", *expectedAuthor.Email)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("slug")
		c.SetParamValues(slug)
		ctx := c.Request().Context()
		articleUnpublisherMock.EXPECT().GetArticleBySlug(ctx, slug).Return(nil, app.ArticleNotFoundError(slug, nil)).Once()

		// Act
		err := handler.UnpublishArticle(c)

		// Assert
		require.ErrorContains(t, err, api.ArticleNotFound(slug).Error())
	})

	t.Run("Should only delete aritcles authored by the currently authenticated user", func(t *testing.T) {
		// Arrange
		expectedAuthor := assembleRandomUser()
		expectedArticle := assembleArticleModel(*expectedAuthor.ID)
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/article/%s", *expectedArticle.Slug), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", primitive.NewObjectID().Hex())
		req.Header.Set("Goduit-Client-Username", "not-the-author")
		req.Header.Set("Goduit-Subject", "not.the.author.email@test.test")
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
