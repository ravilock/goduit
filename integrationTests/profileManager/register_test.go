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

	"github.com/ravilock/goduit/api"
	profileManagerModels "github.com/ravilock/goduit/internal/profileManager/models"

	"github.com/labstack/echo/v4"
	_ "github.com/ravilock/goduit/internal/config"

	"github.com/ravilock/goduit/internal/mongo"
	profileManagerRepositories "github.com/ravilock/goduit/internal/profileManager/repositories"
	profileManagerRequests "github.com/ravilock/goduit/internal/profileManager/requests"
	profileManagerResponses "github.com/ravilock/goduit/internal/profileManager/responses"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	// Connect Mongo DB
	client, err := mongo.ConnectDatabase(viper.GetString("db.url"))
	if err != nil {
		log.Fatalln("Error connecting to database", err)
	}
	profileManagerRepository := profileManagerRepositories.NewUserRepository(client)
	serverUrl := viper.GetString("server.url")
	registerEndpoint := fmt.Sprintf("%s%s", serverUrl, "/api/users")
	httpClient := http.Client{}
	t.Run("Should create new user", func(t *testing.T) {
		registerRequest := generateRegisterBody()
		requestBody, err := json.Marshal(registerRequest)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodPost, registerEndpoint, bytes.NewBuffer(requestBody))
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res, err := httpClient.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()
		require.Equal(t, http.StatusCreated, res.StatusCode)
		registerResponse := new(profileManagerResponses.User)
		resBytes, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		err = json.Unmarshal(resBytes, registerResponse)
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
		req, err := http.NewRequest(http.MethodPost, registerEndpoint, bytes.NewBuffer(requestBody))
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res, err := httpClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusConflict, res.StatusCode)
		resBytes, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		httpError := new(echo.HTTPError)
		err = json.Unmarshal(resBytes, httpError)
		require.NoError(t, err)
		require.Equal(t, httpError.Message, api.ConfictError.Message)
	})
	t.Run("Should not create user with duplicated username", func(t *testing.T) {
		registerRequest := generateRegisterBody()
		registerRequest.User.Email = "different-email@test.test"
		requestBody, err := json.Marshal(registerRequest)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodPost, registerEndpoint, bytes.NewBuffer(requestBody))
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res, err := httpClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusConflict, res.StatusCode)
		resBytes, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		httpError := new(echo.HTTPError)
		err = json.Unmarshal(resBytes, httpError)
		require.NoError(t, err)
		require.Equal(t, httpError.Message, api.ConfictError.Message)
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
