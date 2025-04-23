package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api/validators"
	articlePublisherResponses "github.com/ravilock/goduit/internal/articlePublisher/responses"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestFeedArticles(t *testing.T) {
	err := validators.InitValidator()
	require.NoError(t, err)
	articleFeederMock := newMockArticleFeeder(t)
	profileGetterMock := newMockProfileGetter(t)
	handler := &FeedArticlesHandler{articleFeederMock, profileGetterMock}
	e := echo.New()

	t.Run("Should feed all articles", func(t *testing.T) {
		// Arrange
		limit := 30
		userID := primitive.NewObjectID()
		user := assembleArticleAuthor(userID.Hex())
		articleAuthorID := primitive.NewObjectID()
		expectedArticles := assembleRandomArticles(limit, articleAuthorID)
		expectedAuthor := assembleArticleAuthor(articleAuthorID.Hex())
		req := httptest.NewRequest(http.MethodGet, "/api/articles", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", user.ID.Hex())
		req.Header.Set("Goduit-Client-Username", *user.Username)
		req.Header.Set("Goduit-Client-Email", *user.Email)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		urlValues := c.QueryParams()
		urlValues.Add("limit", strconv.Itoa(limit))
		c.Request().URL.RawQuery = urlValues.Encode()
		ctx := c.Request().Context()
		articleFeederMock.EXPECT().FeedArticles(ctx, user.ID.Hex(), int64(limit), int64(0)).Return(expectedArticles, nil).Once()
		profileGetterMock.EXPECT().GetProfileByID(ctx, articleAuthorID.Hex()).Return(expectedAuthor, nil).Times(limit)

		// Act
		err := handler.FeedArticles(c)

		// Assert
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rec.Code)
		feedArticlesResponse := new(articlePublisherResponses.ArticlesResponse)
		err = json.Unmarshal(rec.Body.Bytes(), feedArticlesResponse)
		require.NoError(t, err)
		checkFeedArticlesResponse(t, limit, feedArticlesResponse)
	})
}

func checkFeedArticlesResponse(t *testing.T, limit int, response *articlePublisherResponses.ArticlesResponse) {
	t.Helper()
	require.LessOrEqual(t, len(response.Articles), limit)
}
