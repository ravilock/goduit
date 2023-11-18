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
	"github.com/ravilock/goduit/api/responses"
	"github.com/ravilock/goduit/internal/config/mongo"
	"github.com/ravilock/goduit/internal/profileManager/repositories"
	"github.com/ravilock/goduit/internal/profileManager/requests"
	"github.com/ravilock/goduit/internal/profileManager/services"
	"github.com/stretchr/testify/assert"
)

const oldUpdateProfileTestUsername = "update-profile-test-username"
const updateProfileTestUsername = "update-profile-test-username"
const updateProfileTestEmail = "update.profile.email@test.test"
const updateProfileTestPassword = "update-profile-test-password"
const updateProfileTestBio = "update profile test bio"
const updateProfileTestImage = "https://update.profile.test.bio.com/image"

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
	repository := repositories.NewUserRepository(client)
	manager := services.NewProfileManager(repository)
	handler := NewProfileHandler(manager)
	clearDatabase(client)
	e := echo.New()
	if err := registerUser(oldUpdateProfileTestUsername, updateProfileTestEmail, "", handler.registerProfileHandler); err != nil {
		log.Fatal("Could not create user", err)
	}
	t.Run("Should fully update an authenticated user's profile", func(t *testing.T) {
		updateProfileRequest := generateUpdateProfileBody()
		requestBody, err := json.Marshal(updateProfileRequest)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPut, "/user", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", updateProfileTestEmail)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = handler.UpdateProfile(c)
		assert.NoError(t, err)
		if rec.Code != http.StatusOK {
			t.Errorf("Got status different than %v, got %v", http.StatusOK, rec.Code)
		}
		updateProfileResponse := new(responses.User)
		err = json.Unmarshal(rec.Body.Bytes(), updateProfileResponse)
		assert.NoError(t, err)
		checkUpdateProfileResponse(t, updateProfileRequest, updateProfileResponse)
	})
	t.Run("Should update only the requested fields", func(t *testing.T) {
		request := new(requests.UpdateProfile)
		request.User.Username = oldUpdateProfileTestUsername
		request.User.Email = updateProfileTestEmail
		requestBody, err := json.Marshal(request)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPut, "/user", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", updateProfileTestEmail)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = handler.UpdateProfile(c)
		assert.NoError(t, err)
		if rec.Code != http.StatusOK {
			t.Errorf("Got status different than %v, got %v", http.StatusOK, rec.Code)
		}
		updateProfileResponse := new(responses.User)
		err = json.Unmarshal(rec.Body.Bytes(), updateProfileResponse)
		assert.NoError(t, err)
		assert.Equal(t, request.User.Username, updateProfileResponse.User.Username, "User email should be the same")
		assert.Equal(t, updateProfileTestEmail, updateProfileResponse.User.Email, "User email should be the same")
		assert.Equal(t, "", updateProfileResponse.User.Bio, "User username should be the same")
		assert.Equal(t, "", updateProfileResponse.User.Image, "User username should be the same")
	})
}

func generateUpdateProfileBody() *requests.UpdateProfile {
	request := new(requests.UpdateProfile)
	request.User.Username = updateProfileTestUsername
	request.User.Email = updateProfileTestEmail
	request.User.Password = updateProfileTestPassword
	request.User.Bio = updateProfileTestBio
	request.User.Image = updateProfileTestImage
	return request
}

func checkUpdateProfileResponse(t *testing.T, request *requests.UpdateProfile, response *responses.User) {
	t.Helper()
	assert.Equal(t, request.User.Email, response.User.Email, "User email should be the same")
	assert.Equal(t, request.User.Username, response.User.Username, "User username should be the same")
	assert.Equal(t, request.User.Bio, response.User.Bio, "User username should be the same")
	assert.Equal(t, request.User.Image, response.User.Image, "User username should be the same")
}
