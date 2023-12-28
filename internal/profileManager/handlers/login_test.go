package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/api/responses"
	"github.com/ravilock/goduit/internal/config/mongo"
	followerCentralRepositories "github.com/ravilock/goduit/internal/followerCentral/repositories"
	followerCentral "github.com/ravilock/goduit/internal/followerCentral/services"
	profileManagerRepositories "github.com/ravilock/goduit/internal/profileManager/repositories"
	profileManagerRequests "github.com/ravilock/goduit/internal/profileManager/requests"
	profileManager "github.com/ravilock/goduit/internal/profileManager/services"
	"github.com/stretchr/testify/assert"
)

const loginTestUsername = "login-test-username"
const loginTestEmail = "login.test.email@test.test"
const loginTestPassword = "login-test-password"

func TestLogin(t *testing.T) {
	databaseURI := os.Getenv("DB_URL")
	if databaseURI == "" {
		log.Fatalln("You must sey your 'DATABASE_URI' environmental variable.")
	}
	// Connect Mongo DB
	client, err := mongo.ConnectDatabase(databaseURI)
	if err != nil {
		log.Fatalln("Error connecting to database", err)
	}
	followerCentralRepository := followerCentralRepositories.NewFollowerRepository(client)
	followerCentral := followerCentral.NewFollowerCentral(followerCentralRepository)
	profileManagerRepository := profileManagerRepositories.NewUserRepository(client)
	profileManager := profileManager.NewProfileManager(profileManagerRepository)
	handler := NewProfileHandler(profileManager, followerCentral)
	clearDatabase(client)
	e := echo.New()
	if err := registerUser(loginTestUsername, loginTestEmail, loginTestPassword, handler.registerProfileHandler); err != nil {
		log.Fatal("Could not create user", err)
	}
	t.Run("Should successfully login", func(t *testing.T) {
		loginRequest := generateLoginBody()
		requestBody, err := json.Marshal(loginRequest)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = handler.Login(c)
		assert.NoError(t, err)
		if rec.Code != http.StatusOK {
			t.Errorf("Got status different than %v, got %v", http.StatusOK, rec.Code)
		}
		loginResponse := new(responses.User)
		err = json.Unmarshal(rec.Body.Bytes(), loginResponse)
		assert.NoError(t, err)
		checkLoginResponse(t, loginRequest, loginResponse)
	})
	t.Run("Should return 401 if email is not found", func(t *testing.T) {
		loginRequest := generateLoginBody()
		loginRequest.User.Email = "wrong-email@test.test"
		requestBody, err := json.Marshal(loginRequest)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = handler.Login(c)
		assert.ErrorIs(t, err, api.FailedLoginAttempt)
	})
	t.Run("Should return 401 if wrong password", func(t *testing.T) {
		loginRequest := generateLoginBody()
		loginRequest.User.Password = "wrong-user-password"
		requestBody, err := json.Marshal(loginRequest)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = handler.Login(c)
		assert.ErrorIs(t, err, api.FailedLoginAttempt)
	})
}

func generateLoginBody() *profileManagerRequests.Login {
	request := new(profileManagerRequests.Login)
	request.User.Email = loginTestEmail
	request.User.Password = loginTestPassword
	return request
}

func checkLoginResponse(t *testing.T, request *profileManagerRequests.Login, response *responses.User) {
	t.Helper()
	assert.Equal(t, request.User.Email, response.User.Email, "User email should be the same")
	assert.NotZero(t, response.User.Token)
	assert.Zero(t, response.User.Image)
	assert.Zero(t, response.User.Bio)
}
