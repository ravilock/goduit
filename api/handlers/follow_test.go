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
	"github.com/ravilock/goduit/internal/app/models"
	"github.com/ravilock/goduit/internal/app/repositories"
	"github.com/stretchr/testify/assert"
)

func TestFollow(t *testing.T) {
	const followTestUsername = "follow-test-username"
	const followTestEmail = "follow.email@test.test"

	const followerUsername = "follower-username"
	const followerEmail = "follower.email@test.test"

	clearDatabase()
	if err := registerUser(followTestUsername, getUserTestEmail, ""); err != nil {
		log.Fatal("Could not create user", err)
	}
	if err := registerUser(followerUsername, followerEmail, ""); err != nil {
		log.Fatal("Could not create user", err)
	}

	e := echo.New()
	t.Run("Should follow a user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/%s/follow", followTestUsername), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("username")
		c.SetParamValues(followTestUsername)
		req.Header.Set("Goduit-Client-Username", followerUsername)
		err := Follow(c)
		assert.NoError(t, err)
		if rec.Code != http.StatusOK {
			t.Errorf("Got status different than %v, got %v", http.StatusOK, rec.Code)
		}
		followResponse := new(responses.ProfileResponse)
		err = json.Unmarshal(rec.Body.Bytes(), followResponse)
		assert.NoError(t, err)
		checkFollowResponse(t, followTestUsername, true, followResponse)
		followerModel, err := repositories.IsFollowedBy(followTestUsername, followerUsername, context.Background())
		assert.NoError(t, err)
		checkFollowerModel(t, followTestUsername, followerUsername, followerModel)
	})
	t.Run("Should return 404 if no user is found", func(t *testing.T) {
		inexistentUsername := "inexistent-username"
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/%s/follow", inexistentUsername), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("username")
		c.SetParamValues(inexistentUsername)
		req.Header.Set("Goduit-Client-Username", followerUsername)
		err := Follow(c)
		assert.ErrorContains(t, err, api.UserNotFound(inexistentUsername).Error())
	})
}

func checkFollowResponse(t *testing.T, username string, following bool, response *responses.ProfileResponse) {
	t.Helper()
	assert.Equal(t, username, response.Profile.Username, "User username should be the same")
	assert.Equal(t, following, response.Profile.Following)
	assert.Zero(t, response.Profile.Image)
	assert.Zero(t, response.Profile.Bio)
}

func checkFollowerModel(t *testing.T, followed, follower string, model *models.Follower) {
	t.Helper()
	assert.Equal(t, followed, *model.Followed, "Wrong followed username")
	assert.Equal(t, follower, *model.Follower, "Wrong follower username")
}

func followUser(followed, follower string) error {
	return repositories.Follow(followed, follower, context.Background())
}
