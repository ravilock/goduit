package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	articlePublisherRepositories "github.com/ravilock/goduit/internal/articlePublisher/repositories"
	articlePublisherRequests "github.com/ravilock/goduit/internal/articlePublisher/requests"
	articlePublisherResponses "github.com/ravilock/goduit/internal/articlePublisher/responses"
	articlePublisher "github.com/ravilock/goduit/internal/articlePublisher/services"
	"github.com/ravilock/goduit/internal/config/mongo"
	followerCentralRepositories "github.com/ravilock/goduit/internal/followerCentral/repositories"
	followerCentral "github.com/ravilock/goduit/internal/followerCentral/services"
	profileManagerModels "github.com/ravilock/goduit/internal/profileManager/models"
	profileManagerRepositories "github.com/ravilock/goduit/internal/profileManager/repositories"
	profileManager "github.com/ravilock/goduit/internal/profileManager/services"
	"github.com/stretchr/testify/assert"
)

func TestWriteArticle(t *testing.T) {
	const createArticleTestUsername = "create-article-test-username"
	const createArticleTestEmail = "create.article.test@test.test"

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
	_, err = registerUser(createArticleTestUsername, createArticleTestEmail, "", profileManager)
	if err != nil {
		t.Error("Could not create user", err)
	}
	e := echo.New()
	t.Run("Should create an article", func(t *testing.T) {
		createArticleRequest := generateWriteArticleBody()
		requestBody, err := json.Marshal(createArticleRequest)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/api/articles", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Client-Username", createArticleTestUsername)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = handler.WriteArticle(c)
		assert.NoError(t, err)
		if rec.Code != http.StatusCreated {
			t.Errorf("Got status different than %v, got %v", http.StatusCreated, rec.Code)
		}
		createArticleResponse := new(articlePublisherResponses.Article)
		err = json.Unmarshal(rec.Body.Bytes(), createArticleResponse)
		checkWriteArticleResponse(t, createArticleRequest, createArticleTestUsername, createArticleResponse)
		assert.NoError(t, err)
	})
	// TODO: Add test for articles with the same Title/Slug
}

func generateWriteArticleBody() *articlePublisherRequests.WriteArticle {
	request := new(articlePublisherRequests.WriteArticle)
	request.Article.Title = "Test Article Name"
	request.Article.Description = "Test Article Description"
	request.Article.Body = "Test Article Body"
	request.Article.TagList = []string{"test"}
	return request
}

func checkWriteArticleResponse(t *testing.T, request *articlePublisherRequests.WriteArticle, author string, response *articlePublisherResponses.Article) {
	t.Helper()
	assert.Equal(t, request.Article.Title, response.Article.Title, "Wrong article title")
	assert.Equal(t, request.Article.Description, response.Article.Description, "Wrong article description")
	assert.Equal(t, request.Article.Body, response.Article.Body, "Wrong article body")
	assert.Equal(t, makeSlug(request.Article.Title), response.Article.Slug, "Wrong article body")
	assert.Equal(t, request.Article.TagList, response.Article.TagList, "Wrong article body")
	assert.Equal(t, author, response.Article.Author.Username, "Wrong article author username")
}

func makeSlug(title string) string {
	loweredTitle := strings.ToLower(title)
	titleWords := strings.Split(loweredTitle, " ")
	return strings.Join(titleWords, "-")
}

func registerUser(username, email, password string, manager *profileManager.ProfileManager) (string, error) {
	if username == "" {
		username = "default-username"
	}
	if email == "" {
		email = "default.email@test.test"
	}
	if password == "" {
		password = "default-password"
	}
	return manager.Register(context.Background(), &profileManagerModels.User{Username: &username, Email: &email}, password)
}

func createArticle(title, description, body, authorUsername string, tagList []string, handler writeArticleHandler) error {
	if title == "" {
		title = "Default Title"
	}
	if description == "" {
		description = "Default Description"
	}
	if body == "" {
		body = "Default Body"
	}
	if len(tagList) == 0 {
		tagList = []string{"default-tag", "test"}
	}
	request := new(articlePublisherRequests.WriteArticle)
	request.Article.Title = title
	request.Article.Description = description
	request.Article.Body = body
	request.Article.TagList = tagList
	requestBody, err := json.Marshal(request)
	if err != nil {
		return err
	}
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/articles", bytes.NewBuffer(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Goduit-Client-Username", authorUsername)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if err := handler.WriteArticle(c); err != nil {
		return err
	}
	return nil
}
