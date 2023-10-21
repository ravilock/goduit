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
	"github.com/ravilock/goduit/api/responses"
	"github.com/ravilock/goduit/internal/app/models"
	"github.com/ravilock/goduit/internal/app/repositories"
	"github.com/stretchr/testify/assert"
)

const followTestUsername = "follow-test-username"
const followTestEmail = "follow.email@test.test"
const followerUsername = "follower-username"
const followerEmail = "follower.email@test.test"

func TestFollow(t *testing.T) {
	clearDatabase()
	e := echo.New()
	if err := registerAccount(followTestUsername, getUserTestEmail, ""); err != nil {
		log.Fatal("Could not create user", err)
	}
	if err := registerAccount(followerUsername, followerEmail, ""); err != nil {
		log.Fatal("Could not create user", err)
	}
	t.Run("Should follow a user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s/follow", followTestUsername), nil)
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
  // TODO: add 404 error test
}

func checkFollowResponse(t *testing.T, username string, following bool, response *responses.ProfileResponse) {
	t.Helper()
	assert.Equal(t, username, response.Profile.Username, "User username should be the same")
	assert.Equal(t, !following, response.Profile.Following)
	assert.Zero(t, response.Profile.Image)
	assert.Zero(t, response.Profile.Bio)
}

func checkFollowerModel(t *testing.T, followed, follower string, model *models.Follower) {
	t.Helper()
	assert.Equal(t, followed, *model.Followed, "Wrong followed username")
	assert.Equal(t, follower, *model.Follower, "Wrong follower username")
}
