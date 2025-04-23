package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/api/validators"
	"github.com/ravilock/goduit/internal/app"

	"github.com/stretchr/testify/require"
)

func TestDeleteComment(t *testing.T) {
	err := validators.InitValidator()
	require.NoError(t, err)
	commentDeleterMock := newMockCommentDeleter(t)
	commentGetterMock := newMockCommentGetter(t)
	articleGetterMock := newMockArticleGetter(t)
	handler := &DeleteCommentHandler{commentDeleterMock, commentGetterMock, articleGetterMock}
	e := echo.New()

	t.Run("Should delete a commentary", func(t *testing.T) {
		// Arrange
		articleAuthorID := primitive.NewObjectID()
		expectedAuthor := assembleArticleAuthor(articleAuthorID.Hex())
		expectedArticle := assembleArticleModel(articleAuthorID)
		expectedComment := assembleCommentModel(articleAuthorID.Hex(), expectedArticle.ID.Hex())
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/article/%s/comments/%s", *expectedArticle.Slug, expectedComment.ID.Hex()), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", expectedAuthor.ID.Hex())
		req.Header.Set("Goduit-Client-Username", *expectedAuthor.Username)
		req.Header.Set("Goduit-Client-Email", *expectedAuthor.Email)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("slug", "id")
		c.SetParamValues(*expectedArticle.Slug, expectedComment.ID.Hex())
		ctx := c.Request().Context()
		articleGetterMock.EXPECT().GetArticleBySlug(ctx, *expectedArticle.Slug).Return(expectedArticle, nil).Once()
		commentGetterMock.EXPECT().GetCommentByID(ctx, expectedComment.ID.Hex()).Return(expectedComment, nil).Once()
		commentDeleterMock.EXPECT().DeleteComment(ctx, expectedComment.ID.Hex()).Return(nil).Once()

		// Act
		err := handler.DeleteComment(c)

		// Assert
		require.NoError(t, err)
		require.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("Only comment author can delete the comment", func(t *testing.T) {
		// Arrange
		articleAuthorID := primitive.NewObjectID()
		expectedArticle := assembleArticleModel(articleAuthorID)
		expectedComment := assembleCommentModel(articleAuthorID.Hex(), expectedArticle.ID.Hex())
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/article/%s/comments/%s", *expectedArticle.Slug, expectedComment.ID.Hex()), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", primitive.NewObjectID().Hex())
		req.Header.Set("Goduit-Client-Username", "not-the-author")
		req.Header.Set("Goduit-Subject", "not.the.author.email@test.test")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("slug", "id")
		c.SetParamValues(*expectedArticle.Slug, expectedComment.ID.Hex())
		ctx := c.Request().Context()
		articleGetterMock.EXPECT().GetArticleBySlug(ctx, *expectedArticle.Slug).Return(expectedArticle, nil).Once()
		commentGetterMock.EXPECT().GetCommentByID(ctx, expectedComment.ID.Hex()).Return(expectedComment, nil).Once()

		// Act
		err := handler.DeleteComment(c)

		// Assert
		require.ErrorContains(t, err, api.Forbidden.Error())
	})

	t.Run("Should return HTTP 404 if no article is found", func(t *testing.T) {
		// Arrange
		articleAuthorID := primitive.NewObjectID()
		expectedAuthor := assembleArticleAuthor(articleAuthorID.Hex())
		expectedArticle := assembleArticleModel(articleAuthorID)
		expectedComment := assembleCommentModel(articleAuthorID.Hex(), expectedArticle.ID.Hex())
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/article/%s/comments/%s", *expectedArticle.Slug, expectedComment.ID.Hex()), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", expectedAuthor.ID.Hex())
		req.Header.Set("Goduit-Client-Username", *expectedAuthor.Username)
		req.Header.Set("Goduit-Client-Email", *expectedAuthor.Email)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("slug", "id")
		ctx := c.Request().Context()
		c.SetParamValues(*expectedArticle.Slug, expectedComment.ID.Hex())
		articleGetterMock.EXPECT().GetArticleBySlug(ctx, *expectedArticle.Slug).Return(nil, app.ArticleNotFoundError(*expectedArticle.Slug, nil)).Once()

		// Act
		err := handler.DeleteComment(c)

		// Assert
		require.ErrorContains(t, err, api.ArticleNotFound(*expectedArticle.Slug).Error())
	})

	t.Run("Should return HTTP 404 if no comment is found", func(t *testing.T) {
		// Arrange
		articleAuthorID := primitive.NewObjectID()
		expectedAuthor := assembleArticleAuthor(articleAuthorID.Hex())
		expectedArticle := assembleArticleModel(articleAuthorID)
		expectedComment := assembleCommentModel(articleAuthorID.Hex(), expectedArticle.ID.Hex())
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/article/%s/comments/%s", expectedArticle.ID.Hex(), expectedComment.ID.Hex()), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", expectedAuthor.ID.Hex())
		req.Header.Set("Goduit-Client-Username", *expectedAuthor.Username)
		req.Header.Set("Goduit-Client-Email", *expectedAuthor.Email)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("slug", "id")
		ctx := c.Request().Context()
		c.SetParamValues(*expectedArticle.Slug, expectedComment.ID.Hex())
		articleGetterMock.EXPECT().GetArticleBySlug(ctx, *expectedArticle.Slug).Return(expectedArticle, nil).Once()
		commentGetterMock.EXPECT().GetCommentByID(ctx, expectedComment.ID.Hex()).Return(nil, app.CommentNotFoundError(expectedComment.ID.Hex(), nil)).Once()

		// Act
		err := handler.DeleteComment(c)

		// Assert
		require.ErrorContains(t, err, api.CommentNotFound(expectedComment.ID.Hex()).Error())
	})
}
