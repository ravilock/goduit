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
	"github.com/ravilock/goduit/api/validators"
	articlePublisherRequests "github.com/ravilock/goduit/internal/articlePublisher/requests"
	encryptionkeys "github.com/ravilock/goduit/internal/config/encryptionKeys"
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

	if err := encryptionkeys.LoadPrivateKey(privateKeyFile); err != nil {
		log.Fatal("Failed to load private key file content", err)
	}

	if err := privateKeyFile.Close(); err != nil {
		log.Fatal("Failed to close private key file", err)
	}

	publicKeyFile, err := os.Open(os.Getenv("PUBLIC_KEY_LOCATION"))
	if err != nil {
		log.Fatal(err)
	}

	if err := encryptionkeys.LoadPublicKey(publicKeyFile); err != nil {
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

func makeSlug(title string) string {
	loweredTitle := strings.ToLower(title)
	titleWords := strings.Split(loweredTitle, " ")
	return strings.Join(titleWords, "-")
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