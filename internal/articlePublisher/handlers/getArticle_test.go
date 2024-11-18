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
	articlePublisherResponses "github.com/ravilock/goduit/internal/articlePublisher/responses"
	profileManagerModels "github.com/ravilock/goduit/internal/profileManager/models"
	"github.com/stretchr/testify/require"
)

func TestGetArticle(t *testing.T) {
	err := validators.InitValidator()
	require.NoError(t, err)
	articleGetterMock := newMockArticleGetter(t)
	profileGetterMock := newMockProfileGetter(t)
	isFollowedCheckerMock := newMockIsFollowedChecker(t)
	handler := &getArticleHandler{service: articleGetterMock, profileManager: profileGetterMock, followerCentral: isFollowedCheckerMock}

	e := echo.New()

	t.Run("Should get an article", func(t *testing.T) {
		// Arrange
		expectedAuthor := assembleRandomUser()
		expectedArticle := assembleArticleModel(*expectedAuthor.ID)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/article/%s", *expectedArticle.Slug), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("slug")
		c.SetParamValues(*expectedArticle.Slug)
		ctx := c.Request().Context()
		articleGetterMock.EXPECT().GetArticleBySlug(ctx, *expectedArticle.Slug).Return(expectedArticle, nil).Once()
		profileGetterMock.EXPECT().GetProfileByID(ctx, *expectedArticle.Author).Return(expectedAuthor, nil).Once()
		isFollowedCheckerMock.EXPECT().IsFollowedBy(ctx, *expectedArticle.Author, "").Return(false).Once()

		// Act
		err := handler.GetArticle(c)

		// Assert
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rec.Code)
		getArticleResponse := new(articlePublisherResponses.ArticleResponse)
		err = json.Unmarshal(rec.Body.Bytes(), getArticleResponse)
		require.NoError(t, err)
		checkGetArticleResponse(t, expectedArticle, expectedAuthor, getArticleResponse)
	})

	t.Run("Should return HTTP 404 if no article is found", func(t *testing.T) {
		// Arrange
		inexistentSlug := "inexistent-slug"
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/article/%s", inexistentSlug), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("slug")
		c.SetParamValues(inexistentSlug)
		ctx := c.Request().Context()
		articleGetterMock.EXPECT().GetArticleBySlug(ctx, inexistentSlug).Return(nil, app.ArticleNotFoundError(inexistentSlug, nil))

		// Act
		err := handler.GetArticle(c)

		// Assert
		require.ErrorContains(t, err, api.ArticleNotFound(inexistentSlug).Error())
	})
}

func checkGetArticleResponse(t *testing.T, expectedArticle *models.Article, expectedAuthor *profileManagerModels.User, response *articlePublisherResponses.ArticleResponse) {
	t.Helper()
	require.Equal(t, *expectedAuthor.Username, response.Article.Author.Username, "Article's author username is wrong")
	require.Equal(t, *expectedArticle.Title, response.Article.Title, "Wrong article title")
	require.Equal(t, *expectedArticle.Slug, response.Article.Slug, "Wrong article slug")
	require.Equal(t, expectedArticle.TagList, response.Article.TagList, "Wrong article tag list")
}
