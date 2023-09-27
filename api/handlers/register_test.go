package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/api/requests"
	"github.com/ravilock/goduit/api/responses"
	"github.com/ravilock/goduit/api/validators"
	"github.com/ravilock/goduit/internal/app/models"
	"github.com/ravilock/goduit/internal/app/repositories"
	encryptionkeys "github.com/ravilock/goduit/internal/config/encryptionKeys"
	"github.com/ravilock/goduit/internal/config/mongo"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	if err := os.Chdir("../.."); err != nil {
		log.Fatalln("Error chanigng directory", err)
	}

	if err := godotenv.Load(".env.test"); err != nil {
		log.Fatalln("No .env file found", err)
	}

	if err := encryptionkeys.LoadKeys(); err != nil {
		log.Fatalln("Failed to read encrpytion keys", err)
	}

	databaseURI := os.Getenv("DB_URL")
	if databaseURI == "" {
		log.Fatal("You must sey your 'DATABASE_URI' environmental variable.")
	}
	// Connect Mongo DB
	if err := mongo.ConnectDatabase(databaseURI); err != nil {
		log.Fatal("Error connecting to database", err)
	}

	clearDatabase()

	// Start Validator
	validators.InitValidator()
}

func clearDatabase() {
	conduitDb := mongo.DatabaseClient.Database("conduit")
	collections, err := conduitDb.ListCollectionNames(context.Background(), bson.D{})
	if err != nil {
		log.Fatal("Could not list collections", err)
	}
	for _, coll := range collections {
		conduitDb.Collection(coll).DeleteMany(context.Background(), bson.D{})
	}
}

func TestRegister(t *testing.T) {
	e := echo.New()
	t.Run("Should create new user", func(t *testing.T) {
		registerRequest := generateRegisterBody()
		requestBody, err := json.Marshal(registerRequest)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = Register(c)
		assert.NoError(t, err)
		if rec.Code != http.StatusCreated {
			t.Errorf("Got status different than %v, got %v", http.StatusCreated, rec.Code)
		}
		registerResponse := new(responses.User)
		err = json.Unmarshal(rec.Body.Bytes(), registerResponse)
		assert.NoError(t, err)
		checkResponse(t, registerRequest, registerResponse)
		userModel, err := repositories.GetUserByEmail(registerRequest.User.Email, context.Background())
		assert.NoError(t, err)
		checkUserModel(t, registerRequest, userModel)
	})
	t.Run("Should not create user with duplicated email", func(t *testing.T) {
		registerRequest := generateRegisterBody()
		registerRequest.User.Username = "different-username"
		requestBody, err := json.Marshal(registerRequest)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = Register(c)
		assert.ErrorIs(t, err, api.ConfictError)
	})
	t.Run("Should not create user with duplicated username", func(t *testing.T) {
		registerRequest := generateRegisterBody()
		registerRequest.User.Email = "different-email@test.test"
		requestBody, err := json.Marshal(registerRequest)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = Register(c)
		assert.ErrorIs(t, err, api.ConfictError)
	})
}

func generateRegisterBody() *requests.Register {
	request := new(requests.Register)
	request.User.Email = "test.test@test.test"
	request.User.Username = "test-username"
	request.User.Password = "test-password"
	return request
}

func checkResponse(t *testing.T, request *requests.Register, response *responses.User) {
	t.Helper()
	assert.Equal(t, request.User.Email, response.User.Email, "User email should be the same")
	assert.Equal(t, request.User.Username, response.User.Username, "User Username should be the same")
	assert.NotZero(t, response.User.Token)
	assert.Zero(t, response.User.Image)
	assert.Zero(t, response.User.Bio)
}

func checkUserModel(t *testing.T, request *requests.Register, user *models.User) {
	t.Helper()
	assert.Equal(t, request.User.Email, *user.Email, "User email should be the same")
	assert.Equal(t, request.User.Username, *user.Username, "User Username should be the same")
	assert.NotZero(t, *user.PasswordHash)
	assert.Zero(t, *user.Image)
	assert.Zero(t, *user.Bio)
}
