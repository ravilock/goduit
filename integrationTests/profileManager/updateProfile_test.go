package profilemanager

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	integrationtests "github.com/ravilock/goduit/integrationTests"
	"github.com/ravilock/goduit/internal/identity"
	"github.com/ravilock/goduit/internal/mongo"
	profileManagerRepositories "github.com/ravilock/goduit/internal/profileManager/repositories"
	profileManagerRequests "github.com/ravilock/goduit/internal/profileManager/requests"
	profileManagerResponses "github.com/ravilock/goduit/internal/profileManager/responses"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestUpdateProfile(t *testing.T) {
	// Connect Mongo DB
	client, err := mongo.ConnectDatabase(viper.GetString("db.url"))
	if err != nil {
		log.Fatalln("Error connecting to database", err)
	}
	profileManagerRepository := profileManagerRepositories.NewUserRepository(client)
	serverUrl := viper.GetString("server.url")
	updateProfileEndpoint := fmt.Sprintf("%s%s", serverUrl, "/api/user")
	httpClient := http.Client{}
	imageServer := mockValidImageURL(t)
	defer imageServer.Close()

	t.Run("Should fully update an authenticated user's profile", func(t *testing.T) {
		// Arrange
		_, cookie := integrationtests.MustRegisterUser(t, profileManagerRequests.RegisterPayload{})
		updateProfileRequest := generateUpdateProfileBody(imageServer.URL)
		requestBody, err := json.Marshal(updateProfileRequest)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodPut, updateProfileEndpoint, bytes.NewBuffer(requestBody))
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.AddCookie(cookie)
		requestTime := time.Now().UTC().Truncate(time.Millisecond)

		// Act
		res, err := httpClient.Do(req)

		// Assert
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)
		updateProfileResponse := new(profileManagerResponses.User)
		resBytes, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		err = json.Unmarshal(resBytes, updateProfileResponse)
		require.NoError(t, err)
		checkUpdateProfileResponse(t, updateProfileRequest, updateProfileResponse)
		checkProfilePassword(t, updateProfileRequest.User.Username, updateProfileRequest.User.Password, profileManagerRepository)
		checkProfileUpdatedAt(t, updateProfileRequest.User.Username, requestTime, profileManagerRepository)
		cookie = integrationtests.CheckCookie(t, res)
		checkUpdatedToken(t, updateProfileRequest, cookie.Value)
	})

	t.Run("Should not update password if not necessary", func(t *testing.T) {
		// Arrange
		oldUpdateProfileTestPassword := uuid.NewString()
		_, cookie := integrationtests.MustRegisterUser(t, profileManagerRequests.RegisterPayload{Password: oldUpdateProfileTestPassword})
		updateProfileRequest := generateUpdateProfileBody(imageServer.URL)
		updateProfileRequest.User.Password = ""
		requestBody, err := json.Marshal(updateProfileRequest)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodPut, updateProfileEndpoint, bytes.NewBuffer(requestBody))
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.AddCookie(cookie)
		requestTime := time.Now().UTC().Truncate(time.Millisecond)

		// Act
		res, err := httpClient.Do(req)

		// Assert
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)
		updateProfileResponse := new(profileManagerResponses.User)
		resBytes, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		err = json.Unmarshal(resBytes, updateProfileResponse)
		require.NoError(t, err)
		checkUpdateProfileResponse(t, updateProfileRequest, updateProfileResponse)
		checkProfilePassword(t, updateProfileRequest.User.Username, oldUpdateProfileTestPassword, profileManagerRepository)
		checkProfileUpdatedAt(t, updateProfileRequest.User.Username, requestTime, profileManagerRepository)
	})

	t.Run("Should not generate new token if not necessary", func(t *testing.T) {
		// Arrange
		oldUserIdentity, cookie := integrationtests.MustRegisterUser(t, profileManagerRequests.RegisterPayload{})
		updateProfileRequest := generateUpdateProfileBody(imageServer.URL)
		updateProfileRequest.User.Password = ""
		updateProfileRequest.User.Username = oldUserIdentity.Username
		updateProfileRequest.User.Email = oldUserIdentity.UserEmail
		requestBody, err := json.Marshal(updateProfileRequest)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodPut, updateProfileEndpoint, bytes.NewBuffer(requestBody))
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.AddCookie(cookie)
		requestTime := time.Now().UTC().Truncate(time.Millisecond)

		// Act
		res, err := httpClient.Do(req)

		// Assert
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)
		updateProfileResponse := new(profileManagerResponses.User)
		resBytes, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		err = json.Unmarshal(resBytes, updateProfileResponse)
		require.NoError(t, err)
		checkUpdateProfileResponse(t, updateProfileRequest, updateProfileResponse)
		checkProfileUpdatedAt(t, updateProfileRequest.User.Username, requestTime, profileManagerRepository)
	})
}

func generateUpdateProfileBody(imageURL string) *profileManagerRequests.UpdateProfileRequest {
	request := new(profileManagerRequests.UpdateProfileRequest)
	request.User.Username = uuid.NewString()
	request.User.Email = fmt.Sprintf("%s@test.test", request.User.Username)
	request.User.Password = uuid.NewString()
	request.User.Bio = uuid.NewString()
	request.User.Image = imageURL
	return request
}

func checkUpdateProfileResponse(t *testing.T, request *profileManagerRequests.UpdateProfileRequest, response *profileManagerResponses.User) {
	t.Helper()
	require.Equal(t, request.User.Username, response.User.Username, "Updated user's username should be %q, got %q", request.User.Username, response.User.Username)
	require.Equal(t, request.User.Email, response.User.Email, "Updated user's email should be %q, got %q", request.User.Email, response.User.Email)
	require.Equal(t, request.User.Bio, response.User.Bio, "Update user's bio should be %q, got %q", request.User.Bio, response.User.Bio)
	require.Equal(t, request.User.Image, response.User.Image, "Update user's image should be %q, got %q", request.User.Image, response.User.Image)
}

func checkUpdatedToken(t *testing.T, request *profileManagerRequests.UpdateProfileRequest, token string) {
	t.Helper()
	identityClaims, err := identity.FromToken(token)
	require.NoError(t, err)
	require.Equal(t, request.User.Username, identityClaims.Username, "Token's username should be %q, got %q", request.User.Username, identityClaims.Username)
	require.Equal(t, request.User.Email, identityClaims.UserEmail, "Token's email should be %q, got %q", request.User.Email, identityClaims.Subject)
}

func checkProfilePassword(t *testing.T, username, password string, repository *profileManagerRepositories.UserRepository) {
	t.Helper()
	user, err := repository.GetUserByUsername(context.Background(), username)
	require.NoError(t, err)
	err = bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(password))
	require.NoError(t, err)
}

func checkProfileUpdatedAt(t *testing.T, username string, requestTime time.Time, repository *profileManagerRepositories.UserRepository) {
	t.Helper()
	user, err := repository.GetUserByUsername(context.Background(), username)
	require.NoError(t, err)
	require.NotNil(t, user.UpdatedAt)
	require.GreaterOrEqual(t, *user.UpdatedAt, requestTime)
}

func mockValidImageURL(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(echo.HeaderContentType, "image/png")
		w.WriteHeader(200)
		_, err := w.Write(nil)
		require.NoError(t, err)
	}))
}
