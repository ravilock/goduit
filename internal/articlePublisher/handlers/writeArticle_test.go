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
	"github.com/ravilock/goduit/api/responses"
	articlePublisherRepositories "github.com/ravilock/goduit/internal/articlePublisher/repositories"
	articlePublisherRequests "github.com/ravilock/goduit/internal/articlePublisher/requests"
	articlePublisher "github.com/ravilock/goduit/internal/articlePublisher/services"
	"github.com/ravilock/goduit/internal/config/mongo"
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
	profileManagerRepository := profileManagerRepositories.NewUserRepository(client)
	profileManager := profileManager.NewProfileManager(profileManagerRepository)
	handler := NewArticlehandler(articlePublisher, profileManager)

	clearDatabase(client)
	_, _, err = registerUser(createArticleTestUsername, createArticleTestEmail, "", profileManager)
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
		createArticleResponse := new(responses.Article)
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

func checkWriteArticleResponse(t *testing.T, request *articlePublisherRequests.WriteArticle, author string, response *responses.Article) {
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

func registerUser(username, email, password string, manager *profileManager.ProfileManager) (*profileManagerModels.User, string, error) {
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
