package handlers

import (
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
	"github.com/ravilock/goduit/internal/articlePublisher/requests"

	articlePublisherResponses "github.com/ravilock/goduit/internal/articlePublisher/responses"
	"github.com/stretchr/testify/require"
)

func TestListComments(t *testing.T) {
	validators.InitValidator()
	commentListerMock := newMockCommentLister(t)
	articleGetterMock := newMockArticleGetter(t)
	profileGetterMock := newMockProfileGetter(t)
	isFollowedCheckerMock := newMockIsFollowedChecker(t)
	handler := &listCommentsHandler{commentListerMock, articleGetterMock, profileGetterMock, isFollowedCheckerMock}

	e := echo.New()

	t.Run("Should list comments", func(t *testing.T) {
		// Arrange
		expectedAuthor := assembleRandomUser()
		expectedArticle := assembleArticleModel(*expectedAuthor.ID)
		comment1 := assembleCommentModel(*expectedAuthor.ID, *expectedArticle.ID, "test comment body")
		comment2 := assembleCommentModel(*expectedAuthor.ID, *expectedArticle.ID, "test comment body 2")
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/article/%s/comments", *expectedArticle.Slug), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("slug")
		c.SetParamValues(*expectedArticle.Slug)
		ctx := c.Request().Context()
		articleGetterMock.EXPECT().GetArticleBySlug(ctx, *expectedArticle.Slug).Return(expectedArticle, nil).Once()
		commentListerMock.EXPECT().ListComments(ctx, expectedArticle.ID.Hex()).Return([]*models.Comment{comment1, comment2}, nil).Once()
		profileGetterMock.EXPECT().GetProfileByID(ctx, *comment1.Author).Return(expectedAuthor, nil).Once()
		isFollowedCheckerMock.EXPECT().IsFollowedBy(ctx, *comment1.Author, "").Return(false).Once()

		// Act
		err := handler.ListComments(c)

		// Assert
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rec.Code)
		listCommentsResponse := new(articlePublisherResponses.CommentsResponse)
		err = json.Unmarshal(rec.Body.Bytes(), listCommentsResponse)
		require.NoError(t, err)
		checkListCommentsResponse(t, listCommentsResponse, 2, comment1, comment2)
	})

	t.Run("Should return HTTP 404 if no article is found", func(t *testing.T) {
		// Arrange
		slug := "inexistent-article"
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/article/%s/comments", slug), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("slug")
		c.SetParamValues(slug)
		ctx := c.Request().Context()
		articleGetterMock.EXPECT().GetArticleBySlug(ctx, slug).Return(nil, app.ArticleNotFoundError(slug, nil)).Once()

		// Act
		err := handler.ListComments(c)

		// Assert
		require.ErrorContains(t, err, api.ArticleNotFound(slug).Error())
	})
}

func generateListCommentsRequest(articleSlug string) *requests.ArticleSlugRequest {
	return &requests.ArticleSlugRequest{
		Slug: articleSlug,
	}
}

func checkListCommentsResponse(t *testing.T, response *articlePublisherResponses.CommentsResponse, length int, comments ...*models.Comment) {
	require.Len(t, response.Comment, length)
	for index, comment := range comments {
		require.Equal(t, comment.ID.Hex(), response.Comment[index].ID)
		require.Equal(t, *comment.Body, response.Comment[index].Body)
	}
}
