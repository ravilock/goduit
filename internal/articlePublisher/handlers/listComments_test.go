package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/api/validators"
	"github.com/ravilock/goduit/internal/app"

	"github.com/ravilock/goduit/internal/articlePublisher/requests"
	"github.com/ravilock/goduit/internal/profileManager/models"

	articlePublisherModels "github.com/ravilock/goduit/internal/articlePublisher/models"
	articlePublisherResponses "github.com/ravilock/goduit/internal/articlePublisher/responses"
	"github.com/stretchr/testify/require"
)

func TestListComments(t *testing.T) {
	err := validators.InitValidator()
	require.NoError(t, err)
	commentListerMock := newMockCommentLister(t)
	articleGetterMock := newMockArticleGetter(t)
	profileGetterMock := newMockProfileGetter(t)
	isFollowedCheckerMock := newMockIsFollowedChecker(t)
	handler := &ListCommentsHandler{commentListerMock, articleGetterMock, profileGetterMock, isFollowedCheckerMock}

	e := echo.New()

	t.Run("Should list comments", func(t *testing.T) {
		// Arrange
		articleAuthorID := primitive.NewObjectID()
		articleAuthor := assembleArticleAuthor(articleAuthorID.Hex())
		authorMap := map[string]*models.User{
			articleAuthorID.Hex(): articleAuthor,
		}
		expectedArticle := assembleArticleModel(articleAuthorID)
		comment1 := assembleCommentModel(articleAuthorID.Hex(), expectedArticle.ID.Hex())
		comment2 := assembleCommentModel(articleAuthorID.Hex(), expectedArticle.ID.Hex())
		listCommentsRequest := generateListCommentsRequest(*expectedArticle.Slug)
		requestBody, err := json.Marshal(listCommentsRequest)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/article/%s/comments", listCommentsRequest.Slug), bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("slug")
		c.SetParamValues(*expectedArticle.Slug)
		ctx := c.Request().Context()
		articleGetterMock.EXPECT().GetArticleBySlug(ctx, *expectedArticle.Slug).Return(expectedArticle, nil).Once()
		commentListerMock.EXPECT().ListComments(ctx, expectedArticle.ID.Hex()).Return([]*articlePublisherModels.Comment{comment1, comment2}, nil).Once()
		profileGetterMock.EXPECT().GetProfileByID(ctx, articleAuthorID.Hex()).Return(articleAuthor, nil).Once()
		isFollowedCheckerMock.EXPECT().IsFollowedBy(ctx, articleAuthorID.Hex(), "").Return(false).Once()

		// Act
		err = handler.ListComments(c)

		// Assert
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rec.Code)
		listCommentsResponse := new(articlePublisherResponses.CommentsResponse)
		err = json.Unmarshal(rec.Body.Bytes(), listCommentsResponse)
		require.NoError(t, err)
		checkListCommentsResponse(t, listCommentsResponse, 2, authorMap, comment1, comment2)
	})

	t.Run("Should return empty array if no comments are found", func(t *testing.T) {
		// Arrange
		articleAuthorID := primitive.NewObjectID()
		articleAuthor := assembleArticleAuthor(articleAuthorID.Hex())
		authorMap := map[string]*models.User{
			articleAuthorID.Hex(): articleAuthor,
		}
		expectedArticle := assembleArticleModel(articleAuthorID)
		listCommentsRequest := generateListCommentsRequest(*expectedArticle.Slug)
		requestBody, err := json.Marshal(listCommentsRequest)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/article/%s/comments", listCommentsRequest.Slug), bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("slug")
		c.SetParamValues(*expectedArticle.Slug)
		ctx := c.Request().Context()
		articleGetterMock.EXPECT().GetArticleBySlug(ctx, *expectedArticle.Slug).Return(expectedArticle, nil).Once()
		commentListerMock.EXPECT().ListComments(ctx, expectedArticle.ID.Hex()).Return([]*articlePublisherModels.Comment{}, nil).Once()

		// Act
		err = handler.ListComments(c)

		// Assert
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rec.Code)
		listCommentsResponse := new(articlePublisherResponses.CommentsResponse)
		err = json.Unmarshal(rec.Body.Bytes(), listCommentsResponse)
		require.NoError(t, err)
		checkListCommentsResponse(t, listCommentsResponse, 0, authorMap)
	})

	t.Run("Should return HTTP 404 if no article is found", func(t *testing.T) {
		// Arrange
		articleSlug := uuid.NewString()
		listCommentsRequest := generateListCommentsRequest(articleSlug)
		requestBody, err := json.Marshal(listCommentsRequest)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/article/%s/comments", listCommentsRequest.Slug), bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("slug")
		c.SetParamValues(articleSlug)
		ctx := c.Request().Context()
		articleGetterMock.EXPECT().GetArticleBySlug(ctx, articleSlug).Return(nil, app.ArticleNotFoundError(articleSlug, nil)).Once()

		// Act
		err = handler.ListComments(c)

		// Assert
		require.ErrorContains(t, err, api.ArticleNotFound(articleSlug).Error())
	})
}

func generateListCommentsRequest(articleSlug string) *requests.ArticleSlugRequest {
	return &requests.ArticleSlugRequest{
		Slug: articleSlug,
	}
}

func checkListCommentsResponse(t *testing.T, response *articlePublisherResponses.CommentsResponse, length int, authorMap map[string]*models.User, comments ...*articlePublisherModels.Comment) {
	require.NotNil(t, response.Comment)
	require.Len(t, response.Comment, length)
	for index, comment := range comments {
		require.Equal(t, comment.ID.Hex(), response.Comment[index].ID)
		require.Equal(t, *comment.Body, response.Comment[index].Body)
		require.Equal(t, comment.CreatedAt, response.Comment[index].CreatedAt)
		require.Equal(t, comment.UpdatedAt, response.Comment[index].UpdatedAt)
		author, ok := authorMap[*comment.Author]
		require.True(t, ok)
		require.Equal(t, *author.Username, response.Comment[index].Author.Username)
	}
}

func assembleCommentModel(commentAuthorID string, articleID string) *articlePublisherModels.Comment {
	commentID := primitive.NewObjectID()
	commentBody := "comment body"
	now := time.Now().UTC().Truncate(time.Millisecond)
	return &articlePublisherModels.Comment{
		ID:        &commentID,
		Author:    &commentAuthorID,
		Article:   &articleID,
		Body:      &commentBody,
		CreatedAt: &now,
		UpdatedAt: &now,
	}
}
