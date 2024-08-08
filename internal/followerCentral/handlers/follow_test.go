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
	followerCentralModels "github.com/ravilock/goduit/internal/followerCentral/models"
	followerCentralRepositories "github.com/ravilock/goduit/internal/followerCentral/repositories"
	followerCentral "github.com/ravilock/goduit/internal/followerCentral/services"
	"github.com/ravilock/goduit/internal/mongo"
	profileManagerRepositories "github.com/ravilock/goduit/internal/profileManager/repositories"
	profileManagerResponses "github.com/ravilock/goduit/internal/profileManager/responses"
	profileManager "github.com/ravilock/goduit/internal/profileManager/services"
	"github.com/stretchr/testify/require"
)

func TestFollow(t *testing.T) {
	const followedTestUsername = "followed-test-username"
	const followedTestEmail = "followed.email@test.test"

	const followerUsername = "follower-username"
	const followerEmail = "follower.email@test.test"

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
	handler := NewFollowerHandler(followerCentral, profileManager)

	clearDatabase(client)
	followedIdentity, err := registerUser(followedTestUsername, followedTestEmail, "", profileManager)
	if err != nil {
		log.Fatalf("Could not create user: %s", err)
	}

	followerIdentity, err := registerUser(followerUsername, followerEmail, "", profileManager)
	if err != nil {
		log.Fatalf("Could not create user: %s", err)
	}

	e := echo.New()
	t.Run("Should follow a user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/%s/follow", followedTestUsername), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("username")
		c.SetParamValues(followedTestUsername)
		req.Header.Set("Goduit-Subject", followerIdentity.Subject)
		req.Header.Set("Goduit-Client-Username", followerIdentity.Username)
		req.Header.Set("Goduit-Client-Email", followerIdentity.UserEmail)
		err := handler.Follow(c)
		require.NoError(t, err)
		if rec.Code != http.StatusOK {
			t.Errorf("Got status different than %v, got %v", http.StatusOK, rec.Code)
		}
		followResponse := new(profileManagerResponses.ProfileResponse)
		err = json.Unmarshal(rec.Body.Bytes(), followResponse)
		require.NoError(t, err)
		checkFollowResponse(t, followedTestUsername, true, followResponse)
		followerModel, err := followerCentralRepository.IsFollowedBy(context.Background(), followedIdentity.Subject, followerIdentity.Subject)
		require.NoError(t, err)
		checkFollowerModel(t, followedIdentity.Subject, followerIdentity.Subject, followerModel)
	})
	t.Run("Should return 404 if no user is found", func(t *testing.T) {
		inexistentUsername := "inexistent-username"
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/%s/follow", inexistentUsername), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("username")
		c.SetParamValues(inexistentUsername)
		req.Header.Set("Goduit-Subject", followerIdentity.Subject)
		req.Header.Set("Goduit-Client-Username", followerIdentity.Username)
		req.Header.Set("Goduit-Client-Email", followerIdentity.UserEmail)
		err := handler.Follow(c)
		require.ErrorContains(t, err, api.UserNotFound(inexistentUsername).Error())
	})
}

func checkFollowResponse(t *testing.T, username string, following bool, response *profileManagerResponses.ProfileResponse) {
	t.Helper()
	require.Equal(t, username, response.Profile.Username, "User username should be the same")
	require.Equal(t, following, response.Profile.Following)
	require.Zero(t, response.Profile.Image)
	require.Zero(t, response.Profile.Bio)
}

func checkFollowerModel(t *testing.T, followed, follower string, model *followerCentralModels.Follower) {
	t.Helper()
	require.NotNil(t, model)
	require.Equal(t, followed, *model.Followed, "Wrong followed username")
	require.Equal(t, follower, *model.Follower, "Wrong follower username")
}
