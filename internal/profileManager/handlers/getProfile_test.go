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
	"github.com/ravilock/goduit/internal/config/mongo"
	followerCentralRepositories "github.com/ravilock/goduit/internal/followerCentral/repositories"
	followerCentral "github.com/ravilock/goduit/internal/followerCentral/services"
	profileManagerRepositories "github.com/ravilock/goduit/internal/profileManager/repositories"
	profileManagerResponses "github.com/ravilock/goduit/internal/profileManager/responses"
	profileManager "github.com/ravilock/goduit/internal/profileManager/services"
	"github.com/stretchr/testify/require"
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
	if _, err := registerUser(getProfileTestUsername, getProfileTestEmail, "", handler.registerProfileHandler); err != nil {
		log.Fatalf("Could not create user")
	}
	t.Run("Should get user profile", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/profiles/%s", getProfileTestUsername), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("username")
		c.SetParamValues(getProfileTestUsername)
		err := handler.GetProfile(c)
		require.NoError(t, err)
		if rec.Code != http.StatusOK {
			t.Errorf("Got status different than %v, got %v", http.StatusOK, rec.Code)
		}
		getProfileResponse := new(profileManagerResponses.ProfileResponse)
		err = json.Unmarshal(rec.Body.Bytes(), getProfileResponse)
		require.NoError(t, err)
		checkGetProfileResponse(t, getProfileTestUsername, false, getProfileResponse)
	})
	t.Run("Should return following as true if logged user follows profile", func(t *testing.T) {
		followerUsername := "follower-username"
		followerEmail := "follower.email@test.test"
		followerIdentity, err := registerUser(followerUsername, followerEmail, "", handler.registerProfileHandler)
		require.NoError(t, err, "Could not register user")
		err = followerCentralRepository.Follow(context.Background(), getProfileTestUsername, followerUsername)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/profiles/%s", getProfileTestUsername), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", followerIdentity.Subject)
		req.Header.Set("Goduit-Client-Username", followerIdentity.Username)
		req.Header.Set("Goduit-Client-Email", followerIdentity.UserEmail)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("username")
		c.SetParamValues(getProfileTestUsername)
		err = handler.GetProfile(c)
		require.NoError(t, err)
		if rec.Code != http.StatusOK {
			t.Errorf("Got status different than %v, got %v", http.StatusOK, rec.Code)
		}
		getProfileResponse := new(profileManagerResponses.ProfileResponse)
		err = json.Unmarshal(rec.Body.Bytes(), getProfileResponse)
		require.NoError(t, err)
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
		require.ErrorContains(t, err, api.UserNotFound(inexistentUsername).Error())
	})
}

func checkGetProfileResponse(t *testing.T, username string, following bool, response *profileManagerResponses.ProfileResponse) {
	t.Helper()
	require.Equal(t, username, response.Profile.Username, "User username should be the same")
	require.Equal(t, following, response.Profile.Following, "Wrong user following")
	require.Zero(t, response.Profile.Image)
	require.Zero(t, response.Profile.Bio)
}
