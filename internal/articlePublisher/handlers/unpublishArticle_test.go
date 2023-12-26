package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
	articlePublisherRepositories "github.com/ravilock/goduit/internal/articlePublisher/repositories"
	articlePublisher "github.com/ravilock/goduit/internal/articlePublisher/services"
	"github.com/ravilock/goduit/internal/config/mongo"
	followerCentralRepositories "github.com/ravilock/goduit/internal/followerCentral/repositories"
	followerCentral "github.com/ravilock/goduit/internal/followerCentral/services"
	profileManagerRepositories "github.com/ravilock/goduit/internal/profileManager/repositories"
	profileManager "github.com/ravilock/goduit/internal/profileManager/services"
	"github.com/stretchr/testify/assert"
)

func TestUnpublishArticle(t *testing.T) {
	const articleAuthorUsername = "article-author-username"
	const articleAuthorEmail = "article.author.email@test.test"

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
	_, err = registerUser(articleAuthorUsername, "", "", profileManager)
	if err != nil {
		t.Error("Could not create user", err)
	}

	e := echo.New()
	t.Run("Should delete an article", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/article/%s", articleSlug), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Client-Username", articleAuthorUsername)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("slug")
		c.SetParamValues(articleSlug)
		err := createArticle(articleTitle, articleDescription, articleBody, articleAuthorUsername, articleTagList, handler.writeArticleHandler)
		assert.NoError(t, err)
		err = handler.UnpublishArticle(c)
		assert.NoError(t, err)
		if rec.Code != http.StatusNoContent {
			t.Errorf("Got status different than %v, got %v", http.StatusNoContent, rec.Code)
		}
	})

	t.Run("Should only delete aritcles authored by authenticated user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/article/%s", articleSlug), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Client-Username", "not-the-author")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("slug")
		c.SetParamValues(articleSlug)
		err := createArticle(articleTitle, articleDescription, articleBody, articleAuthorUsername, articleTagList, handler.writeArticleHandler)
		assert.NoError(t, err)
		err = handler.UnpublishArticle(c)
		assert.ErrorContains(t, err, api.Forbidden.Error())
	})
}
