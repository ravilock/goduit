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
	"github.com/stretchr/testify/assert"
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
	t.Run("Should fully update an authenticated user's profile", func(t *testing.T) {
		oldUpdateProfileTestUsername := uuid.NewString()
		oldUpdateProfileTestEmail := fmt.Sprintf("%s@test.test", oldUpdateProfileTestUsername)
		err := registerUser(oldUpdateProfileTestUsername, oldUpdateProfileTestEmail, "", handler.registerProfileHandler)
		assert.NoError(t, err)
		updateProfileRequest := generateUpdateProfileBody()
		requestBody, err := json.Marshal(updateProfileRequest)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPut, "/user", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", oldUpdateProfileTestEmail)
		req.Header.Set("Goduit-Client-Username", oldUpdateProfileTestUsername)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = handler.UpdateProfile(c)
		assert.NoError(t, err)
		if rec.Code != http.StatusOK {
			t.Errorf("Got status different than %v, got %v", http.StatusOK, rec.Code)
		}
		updateProfileResponse := new(profileManagerResponses.User)
		err = json.Unmarshal(rec.Body.Bytes(), updateProfileResponse)
		assert.NoError(t, err)
		checkUpdateProfileResponse(t, updateProfileRequest, updateProfileResponse)
		checkUpdatedToken(t, updateProfileRequest, updateProfileResponse.User.Token)
		checkProfilePassword(t, updateProfileRequest.User.Username, updateProfileRequest.User.Password, profileManagerRepository)
		clearDatabase(client)
	})
	t.Run("Should not update password if not necessary", func(t *testing.T) {
		oldUpdateProfileTestUsername := uuid.NewString()
		oldUpdateProfileTestEmail := fmt.Sprintf("%s@test.test", oldUpdateProfileTestUsername)
		oldUpdateProfileTestPassword := uuid.NewString()
		err := registerUser(oldUpdateProfileTestUsername, oldUpdateProfileTestEmail, oldUpdateProfileTestPassword, handler.registerProfileHandler)
		assert.NoError(t, err)
		updateProfileRequest := generateUpdateProfileBody()
		updateProfileRequest.User.Password = ""
		requestBody, err := json.Marshal(updateProfileRequest)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPut, "/user", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", oldUpdateProfileTestEmail)
		req.Header.Set("Goduit-Client-Username", oldUpdateProfileTestUsername)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = handler.UpdateProfile(c)
		assert.NoError(t, err)
		if rec.Code != http.StatusOK {
			t.Errorf("Got status different than %v, got %v", http.StatusOK, rec.Code)
		}
		updateProfileResponse := new(profileManagerResponses.User)
		err = json.Unmarshal(rec.Body.Bytes(), updateProfileResponse)
		assert.NoError(t, err)
		checkUpdateProfileResponse(t, updateProfileRequest, updateProfileResponse)
		checkUpdatedToken(t, updateProfileRequest, updateProfileResponse.User.Token)
		checkProfilePassword(t, updateProfileRequest.User.Username, oldUpdateProfileTestPassword, profileManagerRepository)
	})
	t.Run("Should not generate new token if not necessary", func(t *testing.T) {
		oldUpdateProfileTestUsername := uuid.NewString()
		oldUpdateProfileTestEmail := fmt.Sprintf("%s@test.test", oldUpdateProfileTestUsername)
		err := registerUser(oldUpdateProfileTestUsername, oldUpdateProfileTestEmail, "", handler.registerProfileHandler)
		assert.NoError(t, err)
		updateProfileRequest := generateUpdateProfileBody()
		updateProfileRequest.User.Password = ""
		updateProfileRequest.User.Username = oldUpdateProfileTestUsername
		updateProfileRequest.User.Email = oldUpdateProfileTestEmail
		requestBody, err := json.Marshal(updateProfileRequest)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPut, "/user", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", oldUpdateProfileTestEmail)
		req.Header.Set("Goduit-Client-Username", oldUpdateProfileTestUsername)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = handler.UpdateProfile(c)
		assert.NoError(t, err)
		if rec.Code != http.StatusOK {
			t.Errorf("Got status different than %v, got %v", http.StatusOK, rec.Code)
		}
		updateProfileResponse := new(profileManagerResponses.User)
		err = json.Unmarshal(rec.Body.Bytes(), updateProfileResponse)
		assert.NoError(t, err)
		checkUpdateProfileResponse(t, updateProfileRequest, updateProfileResponse)
		assert.Empty(t, "", "Should have not generated new token")
	})
}

func generateUpdateProfileBody() *profileManagerRequests.UpdateProfileRequest {
	request := new(profileManagerRequests.UpdateProfileRequest)
	request.User.Username = uuid.NewString()
	request.User.Email = fmt.Sprintf("%s@test.test", request.User.Username)
	request.User.Password = uuid.NewString()
	request.User.Bio = uuid.NewString()
	request.User.Image = fmt.Sprintf("https://update.profile.test.image.com/image?img=%s", uuid.NewString())
	return request
}

func checkUpdateProfileResponse(t *testing.T, request *profileManagerRequests.UpdateProfileRequest, response *profileManagerResponses.User) {
	t.Helper()
	assert.Equal(t, request.User.Username, response.User.Username, "Updated user's username should be %q, got %q", request.User.Username, response.User.Username)
	assert.Equal(t, request.User.Email, response.User.Email, "Updated user's email should be %q, got %q", request.User.Email, response.User.Email)
	assert.Equal(t, request.User.Bio, response.User.Bio, "Update user's bio should be %q, got %q", request.User.Bio, response.User.Bio)
	assert.Equal(t, request.User.Image, response.User.Image, "Update user's image should be %q, got %q", request.User.Image, response.User.Image)
}

func checkUpdatedToken(t *testing.T, request *profileManagerRequests.UpdateProfileRequest, token string) {
	t.Helper()
	identityClaims, err := identity.FromToken(token)
	assert.NoError(t, err)
	assert.Equal(t, request.User.Username, identityClaims.Username, "Token's username should be %q, got %q", request.User.Username, identityClaims.Username)
	assert.Equal(t, request.User.Email, identityClaims.Subject, "Token's email should be %q, got %q", request.User.Email, identityClaims.Subject)
}

func checkProfilePassword(t *testing.T, username, password string, repository *profileManagerRepositories.UserRepository) {
	t.Helper()
	user, err := repository.GetUserByUsername(context.Background(), username)
	assert.NoError(t, err)
	err = bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(password))
	assert.NoError(t, err)
}
