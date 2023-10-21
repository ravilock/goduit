package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/api/responses"
	"github.com/ravilock/goduit/internal/app/repositories"
	"github.com/stretchr/testify/assert"
)

const getProfileTestUsername = "get-profile-test-username"
const getProfileTestEmail = "get.profile.email@test.test"

func TestGetProfile(t *testing.T) {
	clearDatabase()
	e := echo.New()
	if err := registerAccount(getProfileTestUsername, getUserTestEmail, ""); err != nil {
		log.Fatal("Could not create user", err)
	}
	t.Run("Should get user profile", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/profiles/%s", getProfileTestUsername), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("username")
		c.SetParamValues(getProfileTestUsername)
		err := GetProfile(c)
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
		err := repositories.Follow(getProfileTestUsername, followerUsername, context.Background())
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/profiles/%s", getProfileTestUsername), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Client-Username", followerUsername)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("username")
		c.SetParamValues(getProfileTestUsername)
		err = GetProfile(c)
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
		err := GetProfile(c)
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
