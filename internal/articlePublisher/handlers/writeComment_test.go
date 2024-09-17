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

func TestWriteComment(t *testing.T) {
	databaseURI := os.Getenv("DB_URL")
	if databaseURI == "" {
		log.Fatalln("You must sey your 'DATABASE_URI' environmental variable.")
	}
	// Connect Mongo DB
	client, err := mongo.ConnectDatabase(databaseURI)
	if err != nil {
		log.Fatalln("Error connecting to database", err)
	}
	commentRepository := articlePublisherRepositories.NewCommentRepository(client)
	articlePublisherRepository := articlePublisherRepositories.NewArticleRepository(client)
	commentPublisher := articlePublisher.NewCommentPublisher(commentRepository)
	articlePublisher := articlePublisher.NewArticlePublisher(articlePublisherRepository)
	followerCentralRepository := followerCentralRepositories.NewFollowerRepository(client)
	followerCentral := followerCentral.NewFollowerCentral(followerCentralRepository)
	profileManagerRepository := profileManagerRepositories.NewUserRepository(client)
	profileManager := profileManager.NewProfileManager(profileManagerRepository)
	handler := NewCommentHandler(commentPublisher, articlePublisher, profileManager, followerCentral)

	clearDatabase(client)
	authorIdentity, err := registerUser("", "", "", profileManager)
	if err != nil {
		log.Fatalf("Could not create user: %s", err)
	}
	article, err := createArticle("", "", "", authorIdentity, []string{}, writeArticleHandler{articlePublisher, profileManager})
	if err != nil {
		log.Fatalf("Could not create article: %s", err)
	}
	e := echo.New()
	t.Run("Should create a commentary", func(t *testing.T) {
		createCommentRequest := generateWriteCommentBody()
		requestBody, err := json.Marshal(createCommentRequest)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/articles/%s/comments", createCommentRequest.Slug), bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", authorIdentity.Subject)
		req.Header.Set("Goduit-Client-Username", authorIdentity.Username)
		req.Header.Set("Goduit-Client-Email", authorIdentity.UserEmail)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("slug")
		c.SetParamValues(article.Article.Slug)
		err = handler.WriteComment(c)
		require.NoError(t, err)
		if rec.Code != http.StatusCreated {
			t.Errorf("Got status different than %v, got %v", http.StatusCreated, rec.Code)
		}
		createCommentResponse := new(articlePublisherResponses.CommentResponse)
		err = json.Unmarshal(rec.Body.Bytes(), createCommentResponse)
		require.NoError(t, err)
		checkWriteCommentResponse(t, createCommentRequest, authorIdentity.Username, createCommentResponse)
	})
	t.Run("Should return HTTP 404 if no article is found", func(t *testing.T) {
		articleSlug := "test-slug"
		createCommentRequest := generateWriteCommentBody()
		requestBody, err := json.Marshal(createCommentRequest)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/articles/%s/comments", createCommentRequest.Slug), bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", authorIdentity.Subject)
		req.Header.Set("Goduit-Client-Username", authorIdentity.Username)
		req.Header.Set("Goduit-Client-Email", authorIdentity.UserEmail)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("slug")
		c.SetParamValues(articleSlug)
		err = handler.WriteComment(c)
		require.ErrorContains(t, err, api.ArticleNotFound(articleSlug).Error())
	})
}

func generateWriteCommentBody() *articlePublisherRequests.WriteCommentRequest {
	request := new(articlePublisherRequests.WriteCommentRequest)
	request.Comment.Body = "Test Comment Body"
	return request
}

func checkWriteCommentResponse(t *testing.T, request *articlePublisherRequests.WriteCommentRequest, author string, response *articlePublisherResponses.CommentResponse) {
	t.Helper()
	require.NotZero(t, response.Comment.ID)
	require.Equal(t, request.Comment.Body, response.Comment.Body, "Wrong comment body")
	require.Equal(t, author, response.Comment.Author.Username, "Wrong comment author username")
}
