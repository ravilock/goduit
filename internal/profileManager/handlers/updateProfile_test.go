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

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/internal/config/mongo"
	followerCentralRepositories "github.com/ravilock/goduit/internal/followerCentral/repositories"
	followerCentral "github.com/ravilock/goduit/internal/followerCentral/services"
	"github.com/ravilock/goduit/internal/identity"
	profileManagerRepositories "github.com/ravilock/goduit/internal/profileManager/repositories"
	profileManagerRequests "github.com/ravilock/goduit/internal/profileManager/requests"
	profileManagerResponses "github.com/ravilock/goduit/internal/profileManager/responses"
	profileManager "github.com/ravilock/goduit/internal/profileManager/services"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestUpdateProfile(t *testing.T) {
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
	imageServer := mockValidImageURL(t)
	defer imageServer.Close()
	t.Run("Should fully update an authenticated user's profile", func(t *testing.T) {
		oldUpdateProfileTestUsername := uuid.NewString()
		oldUpdateProfileTestEmail := fmt.Sprintf("%s@test.test", oldUpdateProfileTestUsername)
		identity, err := registerUser(oldUpdateProfileTestUsername, oldUpdateProfileTestEmail, "", profileManager)
		require.NoError(t, err, "Could Not Create User", err)
		updateProfileRequest := generateUpdateProfileBody(imageServer.URL)
		requestBody, err := json.Marshal(updateProfileRequest)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPut, "/user", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", identity.Subject)
		req.Header.Set("Goduit-Client-Username", identity.Username)
		req.Header.Set("Goduit-Client-Email", identity.UserEmail)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = handler.UpdateProfile(c)
		require.NoError(t, err)
		if rec.Code != http.StatusOK {
			t.Errorf("Got status different than %v, got %v", http.StatusOK, rec.Code)
		}
		updateProfileResponse := new(profileManagerResponses.User)
		err = json.Unmarshal(rec.Body.Bytes(), updateProfileResponse)
		require.NoError(t, err)
		checkUpdateProfileResponse(t, updateProfileRequest, updateProfileResponse)
		checkUpdatedToken(t, updateProfileRequest, updateProfileResponse.User.Token)
		checkProfilePassword(t, updateProfileRequest.User.Username, updateProfileRequest.User.Password, profileManagerRepository)
		clearDatabase(client)
	})
	t.Run("Should not update password if not necessary", func(t *testing.T) {
		oldUpdateProfileTestUsername := uuid.NewString()
		oldUpdateProfileTestEmail := fmt.Sprintf("%s@test.test", oldUpdateProfileTestUsername)
		oldUpdateProfileTestPassword := uuid.NewString()
		identity, err := registerUser(oldUpdateProfileTestUsername, oldUpdateProfileTestEmail, oldUpdateProfileTestPassword, profileManager)
		require.NoError(t, err, "Could Not Create User", err)
		updateProfileRequest := generateUpdateProfileBody(imageServer.URL)
		updateProfileRequest.User.Password = ""
		requestBody, err := json.Marshal(updateProfileRequest)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPut, "/user", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", identity.Subject)
		req.Header.Set("Goduit-Client-Email", identity.UserEmail)
		req.Header.Set("Goduit-Client-Username", identity.Username)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = handler.UpdateProfile(c)
		require.NoError(t, err)
		if rec.Code != http.StatusOK {
			t.Errorf("Got status different than %v, got %v", http.StatusOK, rec.Code)
		}
		updateProfileResponse := new(profileManagerResponses.User)
		err = json.Unmarshal(rec.Body.Bytes(), updateProfileResponse)
		require.NoError(t, err)
		checkUpdateProfileResponse(t, updateProfileRequest, updateProfileResponse)
		checkUpdatedToken(t, updateProfileRequest, updateProfileResponse.User.Token)
		checkProfilePassword(t, updateProfileRequest.User.Username, oldUpdateProfileTestPassword, profileManagerRepository)
	})
	t.Run("Should not generate new token if not necessary", func(t *testing.T) {
		oldUpdateProfileTestUsername := uuid.NewString()
		oldUpdateProfileTestEmail := fmt.Sprintf("%s@test.test", oldUpdateProfileTestUsername)
		identity, err := registerUser(oldUpdateProfileTestUsername, oldUpdateProfileTestEmail, "", profileManager)
		require.NoError(t, err, "Could Not Create User", err)
		updateProfileRequest := generateUpdateProfileBody(imageServer.URL)
		updateProfileRequest.User.Password = ""
		updateProfileRequest.User.Username = oldUpdateProfileTestUsername
		updateProfileRequest.User.Email = oldUpdateProfileTestEmail
		requestBody, err := json.Marshal(updateProfileRequest)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPut, "/user", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", identity.Subject)
		req.Header.Set("Goduit-Client-Username", identity.Username)
		req.Header.Set("Goduit-Client-Email", identity.UserEmail)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = handler.UpdateProfile(c)
		require.NoError(t, err)
		if rec.Code != http.StatusOK {
			t.Errorf("Got status different than %v, got %v", http.StatusOK, rec.Code)
		}
		updateProfileResponse := new(profileManagerResponses.User)
		err = json.Unmarshal(rec.Body.Bytes(), updateProfileResponse)
		require.NoError(t, err)
		checkUpdateProfileResponse(t, updateProfileRequest, updateProfileResponse)
		require.Empty(t, "", "Should have not generated new token")
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

func mockValidImageURL(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "image/png")
		w.WriteHeader(200)
		_, err := w.Write(nil)
		require.NoError(t, err)
	}))
}
