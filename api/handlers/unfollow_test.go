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
	"go.mongodb.org/mongo-driver/mongo"
)

func TestUnfollow(t *testing.T) {
	const unfollowTestUsername = "unfollow-test-username"
	const unfollowTestEmail = "unfollow.email@test.test"

	const followerUsername = "follower-username"
	const followerEmail = "follower.email@test.test"

	clearDatabase()
	if err := registerAccount(unfollowTestUsername, getUserTestEmail, ""); err != nil {
		log.Fatal("Could not create user", err)
	}
	if err := registerAccount(followerUsername, followerEmail, ""); err != nil {
		log.Fatal("Could not create user", err)
	}

	followUser(unfollowTestUsername, followerUsername)

	e := echo.New()
	t.Run("Should unfollow a user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/%s/unfollow", unfollowTestUsername), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("username")
		c.SetParamValues(unfollowTestUsername)
		req.Header.Set("Goduit-Client-Username", followerUsername)
		err := Unfollow(c)
		assert.NoError(t, err)
		if rec.Code != http.StatusOK {
			t.Errorf("Got status different than %v, got %v", http.StatusOK, rec.Code)
		}
		followResponse := new(responses.ProfileResponse)
		err = json.Unmarshal(rec.Body.Bytes(), followResponse)
		assert.NoError(t, err)
		checkFollowResponse(t, unfollowTestUsername, false, followResponse)
		followerModel, err := repositories.IsFollowedBy(unfollowTestUsername, followerUsername, context.Background())
		assert.ErrorIs(t, err, mongo.ErrNoDocuments)
		assert.Nil(t, followerModel)
	})
	t.Run("If the user is already not following the other user, return HTTP 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/%s/unfollow", unfollowTestUsername), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("username")
		c.SetParamValues(unfollowTestUsername)
		req.Header.Set("Goduit-Client-Username", followerUsername)
		err := Unfollow(c)
		assert.NoError(t, err)
		if rec.Code != http.StatusOK {
			t.Errorf("Got status different than %v, got %v", http.StatusOK, rec.Code)
		}
		followResponse := new(responses.ProfileResponse)
		err = json.Unmarshal(rec.Body.Bytes(), followResponse)
		assert.NoError(t, err)
		checkFollowResponse(t, unfollowTestUsername, false, followResponse)
		followerModel, err := repositories.IsFollowedBy(unfollowTestUsername, followerUsername, context.Background())
		assert.ErrorIs(t, err, mongo.ErrNoDocuments)
		assert.Nil(t, followerModel)
	})
	t.Run("Should return 404 if no user is found", func(t *testing.T) {
		inexistentUsername := "inexistent-username"
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/%s/unfollow", inexistentUsername), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("username")
		c.SetParamValues(inexistentUsername)
		req.Header.Set("Goduit-Client-Username", followerUsername)
		err := Unfollow(c)
		assert.ErrorContains(t, err, api.UserNotFound(inexistentUsername).Error())
	})
}
