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
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api/validators"
	articlePublisherRequests "github.com/ravilock/goduit/internal/articlePublisher/requests"
	articlePublisherResponses "github.com/ravilock/goduit/internal/articlePublisher/responses"
	"github.com/ravilock/goduit/internal/config"
	"github.com/ravilock/goduit/internal/identity"
	profileManagerModels "github.com/ravilock/goduit/internal/profileManager/models"
	profileManager "github.com/ravilock/goduit/internal/profileManager/services"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	privateKeyFile, err := os.Open(os.Getenv("PRIVATE_KEY_LOCATION"))
	if err != nil {
		log.Fatal(err)
	}

	if err := config.LoadPrivateKey(privateKeyFile); err != nil {
		log.Fatal("Failed to load private key file content", err)
	}

	if err := privateKeyFile.Close(); err != nil {
		log.Fatal("Failed to close private key file", err)
	}

	publicKeyFile, err := os.Open(os.Getenv("PUBLIC_KEY_LOCATION"))
	if err != nil {
		log.Fatal(err)
	}

	if err := config.LoadPublicKey(publicKeyFile); err != nil {
		log.Fatal("Failed to load public key file content", err)
	}

	if err := publicKeyFile.Close(); err != nil {
		log.Fatal("Failed to close publicKeyFile key file", err)
	}

	// Start Validator
	if err := validators.InitValidator(); err != nil {
		log.Fatalln("Failed to load validator", err)
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
