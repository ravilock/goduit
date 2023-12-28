package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/api/responses"
	articlePublisherRepositories "github.com/ravilock/goduit/internal/articlePublisher/repositories"
	articlePublisher "github.com/ravilock/goduit/internal/articlePublisher/services"
	"github.com/ravilock/goduit/internal/config/mongo"
	followerCentralRepositories "github.com/ravilock/goduit/internal/followerCentral/repositories"
	followerCentral "github.com/ravilock/goduit/internal/followerCentral/services"
	profileManagerRepositories "github.com/ravilock/goduit/internal/profileManager/repositories"
	profileManager "github.com/ravilock/goduit/internal/profileManager/services"
	"github.com/stretchr/testify/assert"
)

func TestGetArticle(t *testing.T) {
	const articleAuthorUsername = "article-author-username"

	const articleTitle = "Article Title"
	const articleSlug = "article-title"
	const articleDescription = "Article Description"
	const articleBody = "Article Body"
	articleTagList := []string{"test"}

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
	handler := NewArticlehandler(articlePublisher, profileManager, followerCentral)

	clearDatabase(client)
	_, _, err = registerUser(articleAuthorUsername, "", "", profileManager)
	if err != nil {
		t.Error("Could not create user", err)
	}

	if err := createArticle(articleTitle, articleDescription, articleBody, articleAuthorUsername, articleTagList, handler.writeArticleHandler); err != nil {
		log.Fatal("Could not create article")
	}
	e := echo.New()
	t.Run("Should get an article", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/article/%s", articleSlug), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("slug")
		c.SetParamValues(articleSlug)
		err := handler.GetArticle(c)
		assert.NoError(t, err)
		if rec.Code != http.StatusOK {
			t.Errorf("Got status different than %v, got %v", http.StatusOK, rec.Code)
		}
		getArticleResponse := new(responses.Article)
		err = json.Unmarshal(rec.Body.Bytes(), getArticleResponse)
		assert.NoError(t, err)
		checkGetArticleResponse(t, articleTitle, articleSlug, articleAuthorUsername, articleTagList, getArticleResponse)
	})
	t.Run("Should return http 404 if no article is found", func(t *testing.T) {
		inexistentSlug := "inexistent-slug"
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/article/%s", inexistentSlug), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("slug")
		c.SetParamValues(inexistentSlug)
		err := handler.GetArticle(c)
		assert.ErrorContains(t, err, api.ArticleNotFound(inexistentSlug).Error())
	})
	// TODO: Add test for when the user favorited the article
}

func checkGetArticleResponse(t *testing.T, title, slug, authorUsername string, tagList []string, response *responses.Article) {
	t.Helper()
	assert.Equal(t, authorUsername, response.Article.Author.Username, "Article's author username is wrong")
	assert.Equal(t, title, response.Article.Title, "Wrong article title")
	assert.Equal(t, slug, response.Article.Slug, "Wrong article slug")
	assert.Equal(t, tagList, response.Article.TagList, "Wrong article tag list")
}
