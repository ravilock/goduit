package profilemanager

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
	integrationtests "github.com/ravilock/goduit/integrationTests"
	"github.com/ravilock/goduit/internal/mongo"
	profileManagerRepositories "github.com/ravilock/goduit/internal/profileManager/repositories"
	profileManagerRequests "github.com/ravilock/goduit/internal/profileManager/requests"
	profileManagerResponses "github.com/ravilock/goduit/internal/profileManager/responses"
	"github.com/stretchr/testify/require"

	"github.com/spf13/viper"
)

const (
	loginTestUsername = "login-test-username"
	loginTestEmail    = "login.test.email@test.test"
	loginTestPassword = "login-test-password"
)

func TestLogin(t *testing.T) {
	// Connect Mongo DB
	client, err := mongo.ConnectDatabase(viper.GetString("db.url"))
	if err != nil {
		log.Fatalln("Error connecting to database", err)
	}
	profileManagerRepository := profileManagerRepositories.NewUserRepository(client)
	serverUrl := viper.GetString("server.url")
	loginEndpoint := fmt.Sprintf("%s%s", serverUrl, "/api/users/login")
	httpClient := http.Client{}

	t.Run("Should successfully login", func(t *testing.T) {
		// Arrange
		integrationtests.MustRegisterUser(t, profileManagerRequests.RegisterPayload{Username: loginTestUsername, Email: loginTestEmail, Password: loginTestPassword})
		loginRequest := generateLoginBody()
		preLoginModel, err := profileManagerRepository.GetUserByEmail(context.Background(), loginRequest.User.Email)
		require.NoError(t, err)
		requestBody, err := json.Marshal(loginRequest)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodPost, loginEndpoint, bytes.NewBuffer(requestBody))
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		// Act
		res, err := httpClient.Do(req)

		// Assert
		require.NoError(t, err)
		defer res.Body.Close()
		require.Equal(t, http.StatusOK, res.StatusCode)
		loginResponse := new(profileManagerResponses.User)
		resBytes, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		err = json.Unmarshal(resBytes, loginResponse)
		require.NoError(t, err)
		checkLoginResponse(t, loginRequest, loginResponse)
		postLoginModel, err := profileManagerRepository.GetUserByEmail(context.Background(), loginRequest.User.Email)
		require.NoError(t, err)
		require.GreaterOrEqual(t, *postLoginModel.LastSession, *preLoginModel.LastSession, "User Last Session Was not Updated")
		integrationtests.CheckCookie(t, res)
	})

	t.Run("Should return 401 if email is not found", func(t *testing.T) {
		// Arrange
		loginRequest := generateLoginBody()
		loginRequest.User.Email = "wrong-email@test.test"
		requestBody, err := json.Marshal(loginRequest)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodPost, loginEndpoint, bytes.NewBuffer(requestBody))
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		// Act
		res, err := httpClient.Do(req)

		// Assert
		require.NoError(t, err)
		defer res.Body.Close()
		require.Equal(t, http.StatusUnauthorized, res.StatusCode)
		resBytes, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		httpError := new(echo.HTTPError)
		err = json.Unmarshal(resBytes, httpError)
		require.NoError(t, err)
		require.Equal(t, httpError.Message, api.FailedLoginAttempt.Message)
	})

	t.Run("Should return 401 if password is wrong", func(t *testing.T) {
		// Arrange
		loginRequest := generateLoginBody()
		loginRequest.User.Password = "wrong-user-password"
		requestBody, err := json.Marshal(loginRequest)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodPost, loginEndpoint, bytes.NewBuffer(requestBody))
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		// Act
		res, err := httpClient.Do(req)

		// Assert
		require.NoError(t, err)
		defer res.Body.Close()
		require.Equal(t, http.StatusUnauthorized, res.StatusCode)
		resBytes, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		httpError := new(echo.HTTPError)
		err = json.Unmarshal(resBytes, httpError)
		require.NoError(t, err)
		require.Equal(t, httpError.Message, api.FailedLoginAttempt.Message)
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
	require.Zero(t, response.User.Image)
	require.Zero(t, response.User.Bio)
}
