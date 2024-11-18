package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/internal/articlePublisher/models"
	articlePublisherRequests "github.com/ravilock/goduit/internal/articlePublisher/requests"
	articlePublisherResponses "github.com/ravilock/goduit/internal/articlePublisher/responses"
	"github.com/ravilock/goduit/internal/identity"
	profileManagerModels "github.com/ravilock/goduit/internal/profileManager/models"
	profileManager "github.com/ravilock/goduit/internal/profileManager/services"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func assembleRandomUser() *profileManagerModels.User {
	ID := primitive.NewObjectID()
	username := uuid.NewString()
	email := username + "@test.test"
	image := "http://" + username + ".com"
	now := time.Now()
	return &profileManagerModels.User{
		ID:           &ID,
		Username:     &username,
		Email:        &email,
		PasswordHash: &username,
		Bio:          &username,
		Image:        &image,
		CreatedAt:    &now,
		UpdatedAt:    &now,
		LastSession:  &now,
	}
}

func assembleArticleModel(authorID primitive.ObjectID) *models.Article {
	articleID := primitive.NewObjectID()
	authorIDHex := authorID.Hex()
	articleTitle := "Article Title"
	articleSlug := "article-title"
	articleDescription := "Article Description"
	articleBody := "Article Body"
	articleTagList := []string{"test"}
	favoriteCounts := int64(0)
	now := time.Now()
	return &models.Article{
		ID:             &articleID,
		Author:         &authorIDHex,
		Slug:           &articleSlug,
		Title:          &articleTitle,
		Description:    &articleDescription,
		Body:           &articleBody,
		TagList:        articleTagList,
		CreatedAt:      &now,
		UpdatedAt:      &now,
		FavoritesCount: &favoriteCounts,
	}
}

func assembleCommentModel(authorID, articleID primitive.ObjectID, commentBody string) *models.Comment {
	commentID := primitive.NewObjectID()
	authorIDHex := authorID.Hex()
	articleIDHex := articleID.Hex()
	now := time.Now().UTC().Truncate(time.Millisecond)
	return &models.Comment{
		ID:        &commentID,
		Author:    &authorIDHex,
		Article:   &articleIDHex,
		Body:      &commentBody,
		CreatedAt: &now,
		UpdatedAt: &now,
	}
}

func clearDatabase(client *mongo.Client) {
	conduitDb := client.Database("conduit")
	collections, err := conduitDb.ListCollectionNames(context.Background(), bson.D{})
	if err != nil {
		log.Fatal("Could not list collections", err)
	}
	for _, coll := range collections {
		_, err := conduitDb.Collection(coll).DeleteMany(context.Background(), bson.D{})
		if err != nil {
			log.Fatal("Could not clear database", err)
		}
	}
}

func registerUser(username, email, password string, manager *profileManager.ProfileManager) (*identity.Identity, error) {
	if username == "" {
		username = uuid.NewString()
	}
	if email == "" {
		email = fmt.Sprintf("%s@test.test", uuid.NewString())
	}
	if password == "" {
		password = "default-password"
	}
	token, err := manager.Register(context.Background(), &profileManagerModels.User{Username: &username, Email: &email}, password)
	if err != nil {
		return nil, err
	}
	return identity.FromToken(token)
}

func makeSlug(title string) string {
	loweredTitle := strings.ToLower(title)
	return strings.ReplaceAll(loweredTitle, " ", "-")
}

func createArticles(n int, authorIdentity *identity.Identity, handler writeArticleHandler) ([]*articlePublisherResponses.ArticleResponse, error) {
	articles := []*articlePublisherResponses.ArticleResponse{}
	body := randomString(2500)
	description := randomString(255)
	tags := []string{randomString(10), authorIdentity.Username, authorIdentity.Subject}
	for i := 0; i < n; i++ {
		title := randomString(255)
		response, err := createArticle(title, description, body, authorIdentity, tags, handler)
		if err != nil {
			return nil, err
		}
		articles = append(articles, response)
	}
	return articles, nil
}

func randomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func createArticle(title, description, body string, authorIdentity *identity.Identity, tagList []string, handler writeArticleHandler) (*articlePublisherResponses.ArticleResponse, error) {
	if title == "" {
		title = "Default Title" + uuid.NewString()
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
	request := new(articlePublisherRequests.WriteArticleRequest)
	request.Article.Title = title
	request.Article.Description = description
	request.Article.Body = body
	request.Article.TagList = tagList
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/articles", bytes.NewBuffer(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Goduit-Subject", authorIdentity.Subject)
	req.Header.Set("Goduit-Client-Username", authorIdentity.Username)
	req.Header.Set("Goduit-Client-Email", authorIdentity.UserEmail)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if err := handler.WriteArticle(c); err != nil {
		return nil, err
	}
	response := new(articlePublisherResponses.ArticleResponse)
	if err := json.Unmarshal(rec.Body.Bytes(), response); err != nil {
		return nil, err
	}
	return response, nil
}

func createComment(comment, articleSlug string, authorIdentity *identity.Identity, handler writeCommentHandler) (*articlePublisherResponses.CommentResponse, error) {
	if comment == "" {
		comment = uuid.NewString()
	}
	request := new(articlePublisherRequests.WriteCommentRequest)
	request.Comment.Body = comment
	request.Slug = articleSlug
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/article/%s/comments", articleSlug), bytes.NewBuffer(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Goduit-Subject", authorIdentity.Subject)
	req.Header.Set("Goduit-Client-Username", authorIdentity.Username)
	req.Header.Set("Goduit-Client-Email", authorIdentity.UserEmail)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if err := handler.WriteComment(c); err != nil {
		return nil, err
	}
	response := new(articlePublisherResponses.CommentResponse)
	if err := json.Unmarshal(rec.Body.Bytes(), response); err != nil {
		return nil, err
	}
	return response, nil
}
