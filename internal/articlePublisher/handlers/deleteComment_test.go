package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/ravilock/goduit/api"
	articlePublisherRepositories "github.com/ravilock/goduit/internal/articlePublisher/repositories"

	articlePublisher "github.com/ravilock/goduit/internal/articlePublisher/services"
	followerCentralRepositories "github.com/ravilock/goduit/internal/followerCentral/repositories"
	followerCentral "github.com/ravilock/goduit/internal/followerCentral/services"
	"github.com/ravilock/goduit/internal/mongo"
	profileManagerRepositories "github.com/ravilock/goduit/internal/profileManager/repositories"
	profileManager "github.com/ravilock/goduit/internal/profileManager/services"
	"github.com/stretchr/testify/require"
)

func TestDeleteComment(t *testing.T) {
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
	t.Run("Should delete a commentary", func(t *testing.T) {
		authorIdentity, err := registerUser("", "", "", profileManager)
		require.NoError(t, err)
		article, err := createArticle("", "", "", authorIdentity, []string{}, articleHandler.writeArticleHandler)
		require.NoError(t, err)
		comment, err := createComment("", article.Article.Slug, authorIdentity, handler.writeCommentHandler)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/article/%s/comments/%s", article.Article.Slug, comment.Comment.ID), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", authorIdentity.Subject)
		req.Header.Set("Goduit-Client-Username", authorIdentity.Username)
		req.Header.Set("Goduit-Client-Email", authorIdentity.UserEmail)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("slug", "id")
		c.SetParamValues(article.Article.Slug, comment.Comment.ID)
		err = handler.DeleteComment(c)
		require.NoError(t, err)
		if rec.Code != http.StatusNoContent {
			t.Errorf("Got status different than %v, got %v", http.StatusNoContent, rec.Code)
		}
		require.NoError(t, err)
	})
	t.Run("Only comment author can delete the comment", func(t *testing.T) {
		authorIdentity, err := registerUser("", "", "", profileManager)
		require.NoError(t, err)
		article, err := createArticle("", "", "", authorIdentity, []string{}, articleHandler.writeArticleHandler)
		require.NoError(t, err)
		comment, err := createComment("", article.Article.Slug, authorIdentity, handler.writeCommentHandler)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/article/%s/comments/%s", article.Article.Slug, comment.Comment.ID), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", uuid.NewString())
		req.Header.Set("Goduit-Client-Username", "not-the-author")
		req.Header.Set("Goduit-Client-Email", "not.the.author.email@test.test")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("slug", "id")
		c.SetParamValues(article.Article.Slug, comment.Comment.ID)
		err = handler.DeleteComment(c)
		require.ErrorContains(t, err, api.Forbidden.Error())
	})
	t.Run("Should return http 404 if no article is found", func(t *testing.T) {
		authorIdentity, err := registerUser("", "", "", profileManager)
		require.NoError(t, err)
		articleSlug := uuid.NewString()
		commentID := uuid.NewString()
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/article/%s/comments/%s", articleSlug, commentID), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", authorIdentity.Subject)
		req.Header.Set("Goduit-Client-Username", authorIdentity.Username)
		req.Header.Set("Goduit-Client-Email", authorIdentity.UserEmail)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("slug", "id")
		c.SetParamValues(articleSlug, commentID)
		err = handler.DeleteComment(c)
		require.ErrorContains(t, err, api.ArticleNotFound(articleSlug).Error())
	})
	t.Run("Should return http 404 if no comment is found", func(t *testing.T) {
		authorIdentity, err := registerUser("", "", "", profileManager)
		require.NoError(t, err)
		articleSlug := uuid.NewString()
		commentID := uuid.NewString()
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/article/%s/comments/%s", articleSlug, commentID), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", authorIdentity.Subject)
		req.Header.Set("Goduit-Client-Username", authorIdentity.Username)
		req.Header.Set("Goduit-Client-Email", authorIdentity.UserEmail)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("slug", "id")
		c.SetParamValues(articleSlug, commentID)
		err = handler.DeleteComment(c)
		require.ErrorContains(t, err, api.ArticleNotFound(articleSlug).Error())
	})
}
