package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/api/requests"
	"github.com/ravilock/goduit/api/responses"
	"github.com/stretchr/testify/assert"
)

func TestCreateArticle(t *testing.T) {
	const createArticleTestUsername = "create-article-test-username"
	const createArticleTestEmail = "create.article.test@test.test"

	clearDatabase()
	if err := registerAccount(createArticleTestUsername, createArticleTestEmail, ""); err != nil {
		log.Fatal("Could not create user", err)
	}

	e := echo.New()
	t.Run("Should create an article", func(t *testing.T) {
		createArticleRequest := generateCreateArticleBody()
		requestBody, err := json.Marshal(createArticleRequest)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/api/articles", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Client-Username", createArticleTestUsername)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = CreateArticle(c)
		assert.NoError(t, err)
		if rec.Code != http.StatusCreated {
			t.Errorf("Got status different than %v, got %v", http.StatusCreated, rec.Code)
		}
		createArticleResponse := new(responses.Article)
		err = json.Unmarshal(rec.Body.Bytes(), createArticleResponse)
		checkCreateArticleResponse(t, createArticleRequest, createArticleTestUsername, createArticleResponse)
		assert.NoError(t, err)
	})
	t.Run("Should return http 404 if author's profile is not found", func(t *testing.T) {
		inexistentUsername := "inexistent-username"
		createArticleRequest := generateCreateArticleBody()
		requestBody, err := json.Marshal(createArticleRequest)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/api/articles", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Client-Username", inexistentUsername)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = CreateArticle(c)
		assert.ErrorContains(t, err, api.UserNotFound(inexistentUsername).Error()) // TODO: Should not be 404, should be 401 (Unauthorized)
	})
	// TODO: Add test for articles with the same Title/Slug
}

func generateCreateArticleBody() *requests.CreateArticle {
	request := new(requests.CreateArticle)
	request.Article.Title = "Test Article Name"
	request.Article.Description = "Test Article Description"
	request.Article.Body = "Test Article Body"
	request.Article.TagList = []string{"test"}
	return request
}

func checkCreateArticleResponse(t *testing.T, request *requests.CreateArticle, author string, response *responses.Article) {
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
