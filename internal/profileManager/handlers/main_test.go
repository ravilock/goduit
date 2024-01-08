package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api/validators"
	encryptionkeys "github.com/ravilock/goduit/internal/config/encryptionKeys"
	"github.com/ravilock/goduit/internal/identity"
	profileManagerRequests "github.com/ravilock/goduit/internal/profileManager/requests"
	profileManagerResponses "github.com/ravilock/goduit/internal/profileManager/responses"
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

func registerUser(username, email, password string, handler registerProfileHandler) (*identity.Identity, error) {
	if username == "" {
		username = "default-username"
	}
	if email == "" {
		email = "default.email@test.test"
	}
	if password == "" {
		password = "default-password"
	}
	registerRequest := new(profileManagerRequests.RegisterRequest)
	registerRequest.User.Username = username
	registerRequest.User.Email = email
	registerRequest.User.Password = password
	requestBody, err := json.Marshal(registerRequest)
	if err != nil {
		return nil, err
	}
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if err := handler.Register(c); err != nil {
		return nil, err
	}
	if rec.Code != http.StatusCreated {
		return nil, fmt.Errorf("Got status different than %v, got %v", http.StatusCreated, rec.Code)
	}
	response := new(profileManagerResponses.User)
	if err := json.Unmarshal(rec.Body.Bytes(), response); err != nil {
		return nil, err
	}
	return identity.FromToken(response.User.Token)
}
