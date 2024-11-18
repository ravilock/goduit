package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"strconv"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api/validators"
	"github.com/ravilock/goduit/internal/articlePublisher/models"
	articlePublisherResponses "github.com/ravilock/goduit/internal/articlePublisher/responses"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestListArticles(t *testing.T) {
	validators.InitValidator()
	articleListerMock := newMockArticleLister(t)
	profileGetterMock := newMockProfileGetter(t)
	isFollowedCheckerMock := newMockIsFollowedChecker(t)
	handler := &listArticlesHandler{articleListerMock, profileGetterMock, isFollowedCheckerMock}
	e := echo.New()

	t.Run("Should list all articles", func(t *testing.T) {
		// Arrange
		limit := 30
		expectedAuthor := assembleRandomUser()
		expectedArticles := assembleRandomArticles(limit, *expectedAuthor.ID)
		req := httptest.NewRequest(http.MethodGet, "/api/articles", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		urlValues := c.QueryParams()
		urlValues.Add("limit", strconv.Itoa(limit))
		c.Request().URL.RawQuery = urlValues.Encode()
		ctx := c.Request().Context()
		articleListerMock.EXPECT().ListArticles(ctx, "", "", int64(limit), int64(0)).Return(expectedArticles, nil).Once()
		profileGetterMock.EXPECT().GetProfileByID(ctx, expectedAuthor.ID.Hex()).Return(expectedAuthor, nil).Times(limit)
		isFollowedCheckerMock.EXPECT().IsFollowedBy(ctx, expectedAuthor.ID.Hex(), "").Return(false).Times(limit)

		// Act
		err := handler.ListArticles(c)

		// Assert
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rec.Code)
		listArticlesResponse := new(articlePublisherResponses.ArticlesResponse)
		err = json.Unmarshal(rec.Body.Bytes(), listArticlesResponse)
		require.NoError(t, err)
		checkListArticlesResponse(t, "", "", limit, listArticlesResponse)
	})

	t.Run("Should filter articles based on tag", func(t *testing.T) {
		// Arrange
		limit := 30
		expectedAuthor := assembleRandomUser()
		possibleArticles := assembleRandomArticles(limit, *expectedAuthor.ID)
		tag := possibleArticles[0].TagList[0]
		expectedArticles := slices.DeleteFunc(possibleArticles, func(article *models.Article) bool {
			return !slices.Contains(article.TagList, tag)
		})
		req := httptest.NewRequest(http.MethodGet, "/api/articles", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		urlValues := c.QueryParams()
		urlValues.Add("limit", strconv.Itoa(limit))
		urlValues.Add("tag", tag)
		c.Request().URL.RawQuery = urlValues.Encode()
		ctx := c.Request().Context()
		articleListerMock.EXPECT().ListArticles(ctx, "", tag, int64(limit), int64(0)).Return(expectedArticles, nil).Once()
		profileGetterMock.EXPECT().GetProfileByID(ctx, expectedAuthor.ID.Hex()).Return(expectedAuthor, nil).Times(limit)
		isFollowedCheckerMock.EXPECT().IsFollowedBy(ctx, expectedAuthor.ID.Hex(), "").Return(false).Times(limit)

		// Act
		err := handler.ListArticles(c)

		// Assert
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rec.Code)
		listArticlesResponse := new(articlePublisherResponses.ArticlesResponse)
		err = json.Unmarshal(rec.Body.Bytes(), listArticlesResponse)
		require.NoError(t, err)
		checkListArticlesResponse(t, "", "", limit, listArticlesResponse)
	})

	t.Run("Should filter articles based on author", func(t *testing.T) {
		// Arrange
		limit := 30
		expectedAuthor := assembleRandomUser()
		expectedArticles := assembleRandomArticles(limit, *expectedAuthor.ID)
		req := httptest.NewRequest(http.MethodGet, "/api/articles", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		urlValues := c.QueryParams()
		urlValues.Add("limit", strconv.Itoa(limit))
		urlValues.Add("author", *expectedAuthor.Username)
		c.Request().URL.RawQuery = urlValues.Encode()
		ctx := c.Request().Context()
		profileGetterMock.EXPECT().GetProfileByUsername(ctx, *expectedAuthor.Username).Return(expectedAuthor, nil).Once()
		articleListerMock.EXPECT().ListArticles(ctx, expectedAuthor.ID.Hex(), "", int64(limit), int64(0)).Return(expectedArticles, nil).Once()
		profileGetterMock.EXPECT().GetProfileByID(ctx, expectedAuthor.ID.Hex()).Return(expectedAuthor, nil).Times(limit)
		isFollowedCheckerMock.EXPECT().IsFollowedBy(ctx, expectedAuthor.ID.Hex(), "").Return(false).Times(limit)

		// Act
		err := handler.ListArticles(c)

		// Assert
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rec.Code)
		listArticlesResponse := new(articlePublisherResponses.ArticlesResponse)
		err = json.Unmarshal(rec.Body.Bytes(), listArticlesResponse)
		require.NoError(t, err)
		checkListArticlesResponse(t, "", "", limit, listArticlesResponse)
	})
}

func assembleRandomArticles(amount int, authorId primitive.ObjectID) []*models.Article {
	articles := make([]*models.Article, 0, amount)
	for range amount {
		articles = append(articles, assembleArticleModel(authorId))
	}
	return articles
}

func checkListArticlesResponse(t *testing.T, tag, author string, limit int, response *articlePublisherResponses.ArticlesResponse) {
	t.Helper()
	require.LessOrEqual(t, len(response.Articles), limit)
	if tag == "" && author == "" {
		return
	}
	for _, article := range response.Articles {
		if tag != "" {
			require.Contains(t, article.TagList, tag)
		}
		if author != "" {
			require.Equal(t, article.Author.Username, author)
		}
	}
}

func checkArticlesAreTheSame(articleA, articleB *articlePublisherResponses.Article) bool {
	return articleA.Slug == articleB.Slug
}
