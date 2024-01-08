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

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/internal/config/mongo"
	followerCentralRepositories "github.com/ravilock/goduit/internal/followerCentral/repositories"
	followerCentral "github.com/ravilock/goduit/internal/followerCentral/services"
	profileManagerModels "github.com/ravilock/goduit/internal/profileManager/models"
	profileManagerRepositories "github.com/ravilock/goduit/internal/profileManager/repositories"
	profileManagerRequests "github.com/ravilock/goduit/internal/profileManager/requests"
	profileManagerResponses "github.com/ravilock/goduit/internal/profileManager/responses"
	profileManager "github.com/ravilock/goduit/internal/profileManager/services"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
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
	t.Run("Should create new user", func(t *testing.T) {
		registerRequest := generateRegisterBody()
		requestBody, err := json.Marshal(registerRequest)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = handler.Register(c)
		require.NoError(t, err)
		if rec.Code != http.StatusCreated {
			t.Errorf("Got status different than %v, got %v", http.StatusCreated, rec.Code)
		}
		registerResponse := new(profileManagerResponses.User)
		err = json.Unmarshal(rec.Body.Bytes(), registerResponse)
		require.NoError(t, err)
		checkRegisterResponse(t, registerRequest, registerResponse)
		userModel, err := profileManagerRepository.GetUserByEmail(context.Background(), registerRequest.User.Email)
		require.NoError(t, err)
		checkUserModel(t, registerRequest, userModel)
	})
	t.Run("Should not create user with duplicated email", func(t *testing.T) {
		registerRequest := generateRegisterBody()
		registerRequest.User.Username = "different-username"
		requestBody, err := json.Marshal(registerRequest)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = handler.Register(c)
		require.ErrorIs(t, err, api.ConfictError)
	})
	t.Run("Should not create user with duplicated username", func(t *testing.T) {
		registerRequest := generateRegisterBody()
		registerRequest.User.Email = "different-email@test.test"
		requestBody, err := json.Marshal(registerRequest)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = handler.Register(c)
		require.ErrorIs(t, err, api.ConfictError)
	})
}

func generateRegisterBody() *profileManagerRequests.RegisterRequest {
	request := new(profileManagerRequests.RegisterRequest)
	request.User.Email = "test.test@test.test"
	request.User.Username = "test-username"
	request.User.Password = "test-password"
	return request
}

func checkRegisterResponse(t *testing.T, request *profileManagerRequests.RegisterRequest, response *profileManagerResponses.User) {
	t.Helper()
	require.Equal(t, request.User.Email, response.User.Email, "User email should be the same")
	require.Equal(t, request.User.Username, response.User.Username, "User Username should be the same")
	require.NotZero(t, response.User.Token)
	require.Zero(t, response.User.Image)
	require.Zero(t, response.User.Bio)
}

func checkUserModel(t *testing.T, request *profileManagerRequests.RegisterRequest, user *profileManagerModels.User) {
	t.Helper()
	require.Equal(t, request.User.Email, *user.Email, "User email should be the same")
	require.Equal(t, request.User.Username, *user.Username, "User Username should be the same")
	require.NotZero(t, *user.PasswordHash)
	require.Zero(t, *user.Image)
	require.Zero(t, *user.Bio)
}
