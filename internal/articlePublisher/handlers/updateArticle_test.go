package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
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

func TestUpdateArticle(t *testing.T) {
	const articleTitle = "Update Article Title"
	const articleSlug = "update-article-title"
	const articleDescription = "Update Article Description"
	const articleBody = "Update Article Body"
	articleTagList := []string{"update", "article"}

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
	t.Run("Should update an article", func(t *testing.T) {
		updateArticleRequest := generateUpdateArticleBody()
		requestBody, err := json.Marshal(updateArticleRequest)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/article/%s", articleSlug), bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", authorIdentity.Subject)
		req.Header.Set("Goduit-Client-Username", authorIdentity.Username)
		req.Header.Set("Goduit-Client-Email", authorIdentity.UserEmail)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("slug")
		c.SetParamValues(articleSlug)
		_, err = createArticle(articleTitle, articleDescription, articleBody, authorIdentity, articleTagList, handler.writeArticleHandler)
		require.NoError(t, err)
		err = handler.UpdateArticle(c)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rec.Code)
		updateArticleResponse := new(articlePublisherResponses.ArticleResponse)
		err = json.Unmarshal(rec.Body.Bytes(), updateArticleResponse)
		require.NoError(t, err)
		checkUpdateArticleResponse(t, updateArticleRequest, authorIdentity.Username, updateArticleResponse, articleTagList)
	})
	t.Run("Should return HTTP 404 if no article is found", func(t *testing.T) {
		updateArticleRequest := generateUpdateArticleBody()
		requestBody, err := json.Marshal(updateArticleRequest)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/article/%s", articleSlug), bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", authorIdentity.Subject)
		req.Header.Set("Goduit-Client-Username", authorIdentity.Username)
		req.Header.Set("Goduit-Client-Email", authorIdentity.UserEmail)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("slug")
		c.SetParamValues(articleSlug)
		require.NoError(t, err)
		err = handler.UpdateArticle(c)
		require.ErrorContains(t, err, api.ArticleNotFound(articleSlug).Error())
	})
	t.Run("Should only update articles authored by the currently authenticaed user", func(t *testing.T) {
		updateArticleRequest := generateUpdateArticleBody()
		requestBody, err := json.Marshal(updateArticleRequest)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/article/%s", articleSlug), bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Client-Username", "not-the-author")
		req.Header.Set("Goduit-Subject", "not.the.author.email@test.test")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("slug")
		c.SetParamValues(articleSlug)
		_, err = createArticle(articleTitle, articleDescription, articleBody, authorIdentity, articleTagList, handler.writeArticleHandler)
		require.NoError(t, err)
		err = handler.UpdateArticle(c)
		require.ErrorContains(t, err, api.Forbidden.Error())
	})
}

func generateUpdateArticleBody() *articlePublisherRequests.UpdateArticleRequest {
	request := new(articlePublisherRequests.UpdateArticleRequest)
	request.Article.Title = "New Article Name"
	request.Article.Description = "New Article Description"
	request.Article.Body = "New Article Body"
	return request
}

func checkUpdateArticleResponse(t *testing.T, request *articlePublisherRequests.UpdateArticleRequest, author string, response *articlePublisherResponses.ArticleResponse, tagList []string) {
	t.Helper()
	require.Equal(t, request.Article.Title, response.Article.Title, "Wrong article title")
	require.Equal(t, request.Article.Description, response.Article.Description, "Wrong article description")
	require.Equal(t, request.Article.Body, response.Article.Body, "Wrong article body")
	require.Equal(t, makeSlug(request.Article.Title), response.Article.Slug, "Wrong article body")
	require.Equal(t, tagList, response.Article.TagList, "Wrong article body")
	require.Equal(t, author, response.Article.Author.Username, "Wrong article author username")
}
