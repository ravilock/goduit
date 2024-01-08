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
	"github.com/ravilock/goduit/internal/config/mongo"
	followerCentralRepositories "github.com/ravilock/goduit/internal/followerCentral/repositories"
	followerCentral "github.com/ravilock/goduit/internal/followerCentral/services"
	profileManagerRepositories "github.com/ravilock/goduit/internal/profileManager/repositories"
	profileManagerRequests "github.com/ravilock/goduit/internal/profileManager/requests"
	profileManagerResponses "github.com/ravilock/goduit/internal/profileManager/responses"
	profileManager "github.com/ravilock/goduit/internal/profileManager/services"
	"github.com/stretchr/testify/require"
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
	if _, err := registerUser(loginTestUsername, loginTestEmail, loginTestPassword, profileManager); err != nil {
    log.Fatalf("Could not create user: %s", err)
	}
	t.Run("Should successfully login", func(t *testing.T) {
		loginRequest := generateLoginBody()
		requestBody, err := json.Marshal(loginRequest)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = handler.Login(c)
		require.NoError(t, err)
		if rec.Code != http.StatusOK {
			t.Errorf("Got status different than %v, got %v", http.StatusOK, rec.Code)
		}
		loginResponse := new(profileManagerResponses.User)
		err = json.Unmarshal(rec.Body.Bytes(), loginResponse)
		require.NoError(t, err)
		checkLoginResponse(t, loginRequest, loginResponse)
	})
	t.Run("Should return 401 if email is not found", func(t *testing.T) {
		loginRequest := generateLoginBody()
		loginRequest.User.Email = "wrong-email@test.test"
		requestBody, err := json.Marshal(loginRequest)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = handler.Login(c)
		require.ErrorIs(t, err, api.FailedLoginAttempt)
	})
	t.Run("Should return 401 if wrong password", func(t *testing.T) {
		loginRequest := generateLoginBody()
		loginRequest.User.Password = "wrong-user-password"
		requestBody, err := json.Marshal(loginRequest)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = handler.Login(c)
		require.ErrorIs(t, err, api.FailedLoginAttempt)
	})
}

func generateLoginBody() *profileManagerRequests.LoginRequest {
	request := new(profileManagerRequests.LoginRequest)
	request.User.Email = loginTestEmail
	request.User.Password = loginTestPassword
	return request
}

func checkLoginResponse(t *testing.T, request *profileManagerRequests.LoginRequest, response *profileManagerResponses.User) {
	t.Helper()
	require.Equal(t, request.User.Email, response.User.Email, "User email should be the same")
	require.NotZero(t, response.User.Token)
	require.Zero(t, response.User.Image)
	require.Zero(t, response.User.Bio)
}
