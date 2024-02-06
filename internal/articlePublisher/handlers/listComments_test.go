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

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	// "github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/api"
	articlePublisherRepositories "github.com/ravilock/goduit/internal/articlePublisher/repositories"
	"github.com/ravilock/goduit/internal/articlePublisher/requests"

	// articlePublisherRequests "github.com/ravilock/goduit/internal/articlePublisher/requests"
	articlePublisherResponses "github.com/ravilock/goduit/internal/articlePublisher/responses"
	articlePublisher "github.com/ravilock/goduit/internal/articlePublisher/services"
	"github.com/ravilock/goduit/internal/config/mongo"
	followerCentralRepositories "github.com/ravilock/goduit/internal/followerCentral/repositories"
	followerCentral "github.com/ravilock/goduit/internal/followerCentral/services"
	profileManagerRepositories "github.com/ravilock/goduit/internal/profileManager/repositories"
	profileManager "github.com/ravilock/goduit/internal/profileManager/services"
	"github.com/stretchr/testify/require"
)

func TestListComments(t *testing.T) {
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
	articleHandler := NewArticleHandler(articlePublisher, profileManager, followerCentral)
	handler := NewCommentHandler(commentPublisher, articlePublisher, profileManager, followerCentral)
	clearDatabase(client)
	e := echo.New()
	t.Run("Should list commentaries", func(t *testing.T) {
		authorIdentity, err := registerUser("", "", "", profileManager)
		if err != nil {
			log.Fatalf("Could not create user: %s", err)
		}
		article, err := createArticle("", "", "", authorIdentity, []string{}, articleHandler.writeArticleHandler)
		if err != nil {
			log.Fatalf("Could not create article: %s", err)
		}
		comment1, err := createComment("", article.Article.Slug, authorIdentity, handler.writeCommentHandler)
		require.NoError(t, err)
		comment2, err := createComment("", article.Article.Slug, authorIdentity, handler.writeCommentHandler)
		require.NoError(t, err)
		listCommentsRequest := generateListCommentsRequest(article.Article.Slug)
		requestBody, err := json.Marshal(listCommentsRequest)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/article/%s/comments", listCommentsRequest.Slug), bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", authorIdentity.Subject)
		req.Header.Set("Goduit-Client-Username", authorIdentity.Username)
		req.Header.Set("Goduit-Client-Email", authorIdentity.UserEmail)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("slug")
		c.SetParamValues(article.Article.Slug)
		err = handler.ListComments(c)
		require.NoError(t, err)
		if rec.Code != http.StatusOK {
			t.Errorf("Got status different than %v, got %v", http.StatusCreated, rec.Code)
		}
		listCommentsResponse := new(articlePublisherResponses.CommentsResponse)
		err = json.Unmarshal(rec.Body.Bytes(), listCommentsResponse)
		require.NoError(t, err)
		checkListCommentsResponse(t, listCommentsResponse, 2, comment1, comment2)
	})
	t.Run("Should return http 404 if no article is found", func(t *testing.T) {
		authorIdentity, err := registerUser("", "", "", profileManager)
		if err != nil {
			log.Fatalf("Could not create user: %s", err)
		}
		articleSlug := uuid.NewString()
		listCommentsRequest := generateListCommentsRequest(articleSlug)
		requestBody, err := json.Marshal(listCommentsRequest)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/article/%s/comments", listCommentsRequest.Slug), bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", authorIdentity.Subject)
		req.Header.Set("Goduit-Client-Username", authorIdentity.Username)
		req.Header.Set("Goduit-Client-Email", authorIdentity.UserEmail)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("slug")
		c.SetParamValues(articleSlug)
		err = handler.ListComments(c)
		require.ErrorContains(t, err, api.ArticleNotFound(articleSlug).Error())
	})
}

func generateListCommentsRequest(articleSlug string) *requests.ArticleSlugRequest {
	return &requests.ArticleSlugRequest{
		Slug: articleSlug,
	}
}

func checkListCommentsResponse(t *testing.T, response *articlePublisherResponses.CommentsResponse, length int, comments ...*articlePublisherResponses.CommentResponse) {
	require.Len(t, response.Comment, length)
	for index, comment := range comments {
		require.EqualValues(t, comment.Comment, response.Comment[index])
	}
}
