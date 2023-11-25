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
	mongoConfig "github.com/ravilock/goduit/internal/config/mongo"
	followerCentralRepositories "github.com/ravilock/goduit/internal/followerCentral/repositories"
	followerCentral "github.com/ravilock/goduit/internal/followerCentral/services"
	profileManagerRepositories "github.com/ravilock/goduit/internal/profileManager/repositories"
	profileManagerResponses "github.com/ravilock/goduit/internal/profileManager/responses"
	profileManager "github.com/ravilock/goduit/internal/profileManager/services"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestUnfollow(t *testing.T) {
	const unfollowTestUsername = "unfollow-test-username"
	const unfollowTestEmail = "unfollow.email@test.test"

	const followerUsername = "follower-username"
	const followerEmail = "follower.email@test.test"

	databaseURI := os.Getenv("DB_URL")
	if databaseURI == "" {
		log.Fatalln("You must sey your 'DATABASE_URI' environmental variable.")
	}
	// Connect Mongo DB
	client, err := mongoConfig.ConnectDatabase(databaseURI)
	if err != nil {
		log.Fatalln("Error connecting to database", err)
	}
	followerCentralRepository := followerCentralRepositories.NewFollowerRepository(client)
	followerCentral := followerCentral.NewFollowerCentral(followerCentralRepository)
	profileManagerRepository := profileManagerRepositories.NewUserRepository(client)
	profileManager := profileManager.NewProfileManager(profileManagerRepository)
	handler := NewFollowerHandler(followerCentral, profileManager)

	clearDatabase(client)
	_, _, err = registerUser(unfollowTestUsername, unfollowTestEmail, "", profileManager)
	if err != nil {
		t.Error("Could not create user", err)
	}

	_, _, err = registerUser(followerUsername, followerEmail, "", profileManager)
	if err != nil {
		t.Error("Could not create user", err)
	}

	if err := followUser(unfollowTestUsername, followerUsername, handler.followUserHandler); err != nil {
		t.Error("Failed to setup user following relationship", err)
	}

	e := echo.New()
	t.Run("Should unfollow a user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/%s/unfollow", unfollowTestUsername), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("username")
		c.SetParamValues(unfollowTestUsername)
		req.Header.Set("Goduit-Client-Username", followerUsername)
		err := handler.Unfollow(c)
		assert.NoError(t, err)
		if rec.Code != http.StatusOK {
			t.Errorf("Got status different than %v, got %v", http.StatusOK, rec.Code)
		}
		followResponse := new(profileManagerResponses.ProfileResponse)
		err = json.Unmarshal(rec.Body.Bytes(), followResponse)
		assert.NoError(t, err)
		checkFollowResponse(t, unfollowTestUsername, false, followResponse)
		followerModel, err := followerCentralRepository.IsFollowedBy(context.Background(), unfollowTestUsername, followerUsername)
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
		err := handler.Unfollow(c)
		assert.NoError(t, err)
		if rec.Code != http.StatusOK {
			t.Errorf("Got status different than %v, got %v", http.StatusOK, rec.Code)
		}
		followResponse := new(profileManagerResponses.ProfileResponse)
		err = json.Unmarshal(rec.Body.Bytes(), followResponse)
		assert.NoError(t, err)
		checkFollowResponse(t, unfollowTestUsername, false, followResponse)
		followerModel, err := followerCentralRepository.IsFollowedBy(context.Background(), unfollowTestUsername, followerUsername)
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
		err := handler.Unfollow(c)
		assert.ErrorContains(t, err, api.UserNotFound(inexistentUsername).Error())
	})
}
