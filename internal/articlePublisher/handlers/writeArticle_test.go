package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	articlePublisherRepositories "github.com/ravilock/goduit/internal/articlePublisher/repositories"
	articlePublisherRequests "github.com/ravilock/goduit/internal/articlePublisher/requests"
	articlePublisherResponses "github.com/ravilock/goduit/internal/articlePublisher/responses"
	articlePublisher "github.com/ravilock/goduit/internal/articlePublisher/services"
	followerCentralRepositories "github.com/ravilock/goduit/internal/followerCentral/repositories"
	followerCentral "github.com/ravilock/goduit/internal/followerCentral/services"
	"github.com/ravilock/goduit/internal/mongo"
	profileManagerRepositories "github.com/ravilock/goduit/internal/profileManager/repositories"
	profileManager "github.com/ravilock/goduit/internal/profileManager/services"
	"github.com/stretchr/testify/require"
)

func TestWriteArticle(t *testing.T) {
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
	authorIdentity, err := registerUser("", "", "", profileManager)
	if err != nil {
		log.Fatalf("Could not create user: %s", err)
	}
	e := echo.New()
	t.Run("Should create an article", func(t *testing.T) {
		createArticleRequest := generateWriteArticleBody()
		requestBody, err := json.Marshal(createArticleRequest)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/api/articles", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", authorIdentity.Subject)
		req.Header.Set("Goduit-Client-Username", authorIdentity.Username)
		req.Header.Set("Goduit-Client-Email", authorIdentity.UserEmail)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = handler.WriteArticle(c)
		require.NoError(t, err)
		if rec.Code != http.StatusCreated {
			t.Errorf("Got status different than %v, got %v", http.StatusCreated, rec.Code)
		}
		createArticleResponse := new(articlePublisherResponses.ArticleResponse)
		err = json.Unmarshal(rec.Body.Bytes(), createArticleResponse)
		require.NoError(t, err)
		checkWriteArticleResponse(t, createArticleRequest, authorIdentity.Username, createArticleResponse)
	})
	// TODO: Add test for articles with the same Title/Slug
}

func generateWriteArticleBody() *articlePublisherRequests.WriteArticleRequest {
	request := new(articlePublisherRequests.WriteArticleRequest)
	request.Article.Title = "Test Article Name"
	request.Article.Description = "Test Article Description"
	request.Article.Body = "Test Article Body"
	request.Article.TagList = []string{"test"}
	return request
}

func checkWriteArticleResponse(t *testing.T, request *articlePublisherRequests.WriteArticleRequest, author string, response *articlePublisherResponses.ArticleResponse) {
	t.Helper()
	require.Equal(t, request.Article.Title, response.Article.Title, "Wrong article title")
	require.Equal(t, request.Article.Description, response.Article.Description, "Wrong article description")
	require.Equal(t, request.Article.Body, response.Article.Body, "Wrong article body")
	require.Equal(t, makeSlug(request.Article.Title), response.Article.Slug, "Wrong article body")
	require.Equal(t, request.Article.TagList, response.Article.TagList, "Wrong article body")
	require.Equal(t, author, response.Article.Author.Username, "Wrong article author username")
}
