package handlers

import (
	"context"
	"encoding/json"
	"fmt"
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
	profileManager "github.com/ravilock/goduit/internal/profileManager/services"
	"github.com/stretchr/testify/assert"
)

const getProfileTestUsername = "get-profile-test-username"
const getProfileTestEmail = "get.profile.email@test.test"

func TestGetProfile(t *testing.T) {
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
	if err := registerUser(getProfileTestUsername, getProfileTestEmail, "", handler.registerProfileHandler); err != nil {
		log.Fatal("Could not create user", err)
	}
	t.Run("Should get user profile", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/profiles/%s", getProfileTestUsername), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("username")
		c.SetParamValues(getProfileTestUsername)
		err := handler.GetProfile(c)
		assert.NoError(t, err)
		if rec.Code != http.StatusOK {
			t.Errorf("Got status different than %v, got %v", http.StatusOK, rec.Code)
		}
		getProfileResponse := new(responses.ProfileResponse)
		err = json.Unmarshal(rec.Body.Bytes(), getProfileResponse)
		assert.NoError(t, err)
		checkGetProfileResponse(t, getProfileTestUsername, false, getProfileResponse)
	})
	t.Run("Should return following as true if logged user follows profile", func(t *testing.T) {
		followerUsername := "follower-username"
		err := followerCentralRepository.Follow(context.Background(), getProfileTestUsername, followerUsername)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/profiles/%s", getProfileTestUsername), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Client-Username", followerUsername)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("username")
		c.SetParamValues(getProfileTestUsername)
		err = handler.GetProfile(c)
		assert.NoError(t, err)
		if rec.Code != http.StatusOK {
			t.Errorf("Got status different than %v, got %v", http.StatusOK, rec.Code)
		}
		getProfileResponse := new(responses.ProfileResponse)
		err = json.Unmarshal(rec.Body.Bytes(), getProfileResponse)
		assert.NoError(t, err)
		checkGetProfileResponse(t, getProfileTestUsername, true, getProfileResponse)
	})
	t.Run("Should return http 404 if no user is found", func(t *testing.T) {
		inexistentUsername := "inexistent-username"
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/profiles/%s", inexistentUsername), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("username")
		c.SetParamValues(inexistentUsername)
		err := handler.GetProfile(c)
		assert.ErrorContains(t, err, api.UserNotFound(inexistentUsername).Error())
	})
}

func checkGetProfileResponse(t *testing.T, username string, following bool, response *responses.ProfileResponse) {
	t.Helper()
	assert.Equal(t, username, response.Profile.Username, "User username should be the same")
	assert.Equal(t, following, response.Profile.Following, "Wrong user following")
	assert.Zero(t, response.Profile.Image)
	assert.Zero(t, response.Profile.Bio)
}
