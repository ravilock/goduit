package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/labstack/echo/v4"
	articlePublisherRepositories "github.com/ravilock/goduit/internal/articlePublisher/repositories"
	articlePublisherResponses "github.com/ravilock/goduit/internal/articlePublisher/responses"
	articlePublisher "github.com/ravilock/goduit/internal/articlePublisher/services"
	followerCentralRepositories "github.com/ravilock/goduit/internal/followerCentral/repositories"
	followerCentral "github.com/ravilock/goduit/internal/followerCentral/services"
	"github.com/ravilock/goduit/internal/mongo"
	profileManagerRepositories "github.com/ravilock/goduit/internal/profileManager/repositories"
	profileManager "github.com/ravilock/goduit/internal/profileManager/services"
	"github.com/stretchr/testify/require"
)

func TestListArticles(t *testing.T) {
	const authorUsername1 = "author-username-1"
	const authorEmail1 = "author.email.1@test.test"
	const authorUsername2 = "author-username-2"
	const authorEmail2 = "author.email.2@test.test"
	databaseURI := os.Getenv("DB_URL")
	if databaseURI == "" {
		log.Fatalln("You must sey your 'DATABASE_URI' environmental variable.")
	}
	// Connect Mongo DB
	client, err := mongo.ConnectDatabase(databaseURI)
	if err != nil {
		log.Fatalln("Error connecting to database", err)
	}
	articlePublisherRepository := articlePublisherRepositories.NewArticleRepository(client)
	articlePublisher := articlePublisher.NewArticlePublisher(articlePublisherRepository)
	followerCentralRepository := followerCentralRepositories.NewFollowerRepository(client)
	followerCentral := followerCentral.NewFollowerCentral(followerCentralRepository)
	profileManagerRepository := profileManagerRepositories.NewUserRepository(client)
	profileManager := profileManager.NewProfileManager(profileManagerRepository)
	handler := NewArticleHandler(articlePublisher, profileManager, followerCentral)

	clearDatabase(client)
	authorIdentity1, err := registerUser(authorUsername1, authorEmail1, "", profileManager)
	if err != nil {
		log.Fatalf("Could not create user: %s", err)
	}
	author1Articles, err := createArticles(15, authorIdentity1, handler.writeArticleHandler)
	if err != nil {
		log.Fatalf("Could not create user: %s", err)
	}

	authorIdentity2, err := registerUser(authorUsername2, authorEmail2, "", profileManager)
	if err != nil {
		log.Fatalf("Could not create user: %s", err)
	}
	_, err = createArticles(15, authorIdentity2, handler.writeArticleHandler)
	if err != nil {
		log.Fatalf("Could not create user: %s", err)
	}

	e := echo.New()
	t.Run("Should list all articles", func(t *testing.T) {
		limit := 30
		req := httptest.NewRequest(http.MethodGet, "/api/articles", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.QueryParams().Add("limit", strconv.Itoa(limit))
		err := handler.ListArticles(c)
		require.NoError(t, err)
		if rec.Code != http.StatusOK {
			t.Errorf("Got status different than %v, got %v", http.StatusOK, rec.Code)
		}
		listArticlesResponse := new(articlePublisherResponses.ArticlesResponse)
		err = json.Unmarshal(rec.Body.Bytes(), listArticlesResponse)
		require.NoError(t, err)
		checkListArticlesResponse(t, "", "", limit, listArticlesResponse)
	})
	t.Run("Should filter articles based on tag", func(t *testing.T) {
		tag := author1Articles[0].Article.TagList[0]
		limit := 30
		req := httptest.NewRequest(http.MethodGet, "/api/articles", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		q := c.QueryParams()
		q.Add("limit", strconv.Itoa(limit))
		q.Add("tag", tag)
		err := handler.ListArticles(c)
		require.NoError(t, err)
		if rec.Code != http.StatusOK {
			t.Errorf("Got status different than %v, got %v", http.StatusOK, rec.Code)
		}
		listArticlesResponse := new(articlePublisherResponses.ArticlesResponse)
		err = json.Unmarshal(rec.Body.Bytes(), listArticlesResponse)
		require.NoError(t, err)
		checkListArticlesResponse(t, tag, "", limit, listArticlesResponse)
		require.Len(t, listArticlesResponse.Articles, 15)
	})
	t.Run("Should filter articles based on author", func(t *testing.T) {
		author := authorIdentity2.Username
		limit := 30
		req := httptest.NewRequest(http.MethodGet, "/api/articles", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		q := c.QueryParams()
		q.Add("limit", strconv.Itoa(limit))
		q.Add("author", author)
		err := handler.ListArticles(c)
		require.NoError(t, err)
		if rec.Code != http.StatusOK {
			t.Errorf("Got status different than %v, got %v", http.StatusOK, rec.Code)
		}
		listArticlesResponse := new(articlePublisherResponses.ArticlesResponse)
		err = json.Unmarshal(rec.Body.Bytes(), listArticlesResponse)
		require.NoError(t, err)
		checkListArticlesResponse(t, "", author, limit, listArticlesResponse)
		require.Len(t, listArticlesResponse.Articles, 15)
	})
	t.Run("Should limit number of results", func(t *testing.T) {
		limit := 3
		req := httptest.NewRequest(http.MethodGet, "/api/articles", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.QueryParams().Add("limit", strconv.Itoa(limit))
		err := handler.ListArticles(c)
		require.NoError(t, err)
		if rec.Code != http.StatusOK {
			t.Errorf("Got status different than %v, got %v", http.StatusOK, rec.Code)
		}
		listArticlesResponse := new(articlePublisherResponses.ArticlesResponse)
		err = json.Unmarshal(rec.Body.Bytes(), listArticlesResponse)
		require.NoError(t, err)
		checkListArticlesResponse(t, "", "", limit, listArticlesResponse)
	})
	t.Run("Should properly offset results", func(t *testing.T) {
		offset := 5
		limit := 15
		req1 := httptest.NewRequest(http.MethodGet, "/api/articles", nil)
		req1.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec1 := httptest.NewRecorder()
		c1 := e.NewContext(req1, rec1)
		q1 := c1.QueryParams()
		q1.Add("limit", strconv.Itoa(limit))
		err := handler.ListArticles(c1)
		require.NoError(t, err)
		if rec1.Code != http.StatusOK {
			t.Errorf("Got status different than %v, got %v", http.StatusOK, rec1.Code)
		}
		listArticlesResponse1 := new(articlePublisherResponses.ArticlesResponse)
		err = json.Unmarshal(rec1.Body.Bytes(), listArticlesResponse1)
		require.NoError(t, err)
		checkListArticlesResponse(t, "", "", limit, listArticlesResponse1)

		req2 := httptest.NewRequest(http.MethodGet, "/api/articles", nil)
		req2.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec2 := httptest.NewRecorder()
		c2 := e.NewContext(req2, rec2)
		q2 := c2.QueryParams()
		q2.Add("limit", strconv.Itoa(limit))
		q2.Add("offset", strconv.Itoa(offset))
		err = handler.ListArticles(c2)
		require.NoError(t, err)
		if rec2.Code != http.StatusOK {
			t.Errorf("Got status different than %v, got %v", http.StatusOK, rec1.Code)
		}
		listArticlesResponse2 := new(articlePublisherResponses.ArticlesResponse)
		err = json.Unmarshal(rec2.Body.Bytes(), listArticlesResponse2)
		require.NoError(t, err)
		checkListArticlesResponse(t, "", "", limit, listArticlesResponse2)

		articles1 := listArticlesResponse1.Articles[5:16]
		articles2 := listArticlesResponse2.Articles[:11]
		for i := 0; i < limit-offset; i++ {
			article1 := &articles1[i]
			article2 := &articles2[i]
			require.True(t, checkArticlesAreTheSame(article1, article2))
		}
	})
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
