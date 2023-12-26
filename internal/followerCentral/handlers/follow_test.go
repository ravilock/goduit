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
	followerCentralModels "github.com/ravilock/goduit/internal/followerCentral/models"
	followerCentralRepositories "github.com/ravilock/goduit/internal/followerCentral/repositories"
	followerCentral "github.com/ravilock/goduit/internal/followerCentral/services"
	profileManagerModels "github.com/ravilock/goduit/internal/profileManager/models"
	profileManagerRepositories "github.com/ravilock/goduit/internal/profileManager/repositories"
	profileManagerResponses "github.com/ravilock/goduit/internal/profileManager/responses"
	profileManager "github.com/ravilock/goduit/internal/profileManager/services"
	"github.com/stretchr/testify/assert"
)

func TestFollow(t *testing.T) {
	const followTestUsername = "follow-test-username"
	const followTestEmail = "follow.email@test.test"

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
	_, err = registerUser(followTestUsername, followTestEmail, "", profileManager)
	if err != nil {
		t.Error("Could not create user", err)
	}

	_, err = registerUser(followerUsername, followerEmail, "", profileManager)
	if err != nil {
		t.Error("Could not create user", err)
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
		err := handler.Follow(c)
		assert.NoError(t, err)
		if rec.Code != http.StatusOK {
			t.Errorf("Got status different than %v, got %v", http.StatusOK, rec.Code)
		}
		followResponse := new(profileManagerResponses.ProfileResponse)
		err = json.Unmarshal(rec.Body.Bytes(), followResponse)
		assert.NoError(t, err)
		checkFollowResponse(t, followTestUsername, true, followResponse)
		followerModel, err := followerCentralRepository.IsFollowedBy(context.Background(), followTestUsername, followerUsername)
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
		err := handler.Follow(c)
		assert.ErrorContains(t, err, api.UserNotFound(inexistentUsername).Error())
	})
}

func checkFollowResponse(t *testing.T, username string, following bool, response *profileManagerResponses.ProfileResponse) {
	t.Helper()
	assert.Equal(t, username, response.Profile.Username, "User username should be the same")
	assert.Equal(t, following, response.Profile.Following)
	assert.Zero(t, response.Profile.Image)
	assert.Zero(t, response.Profile.Bio)
}

func checkFollowerModel(t *testing.T, followed, follower string, model *followerCentralModels.Follower) {
	t.Helper()
	assert.NotNil(t, model)
	assert.Equal(t, followed, *model.Followed, "Wrong followed username")
	assert.Equal(t, follower, *model.Follower, "Wrong follower username")
}

func registerUser(username, email, password string, manager *profileManager.ProfileManager) (string, error) {
	if username == "" {
		username = "default-username"
	}
	if email == "" {
		email = "default.email@test.test"
	}
	if password == "" {
		password = "default-password"
	}
	return manager.Register(context.Background(), &profileManagerModels.User{Username: &username, Email: &email}, password)
}

func followUser(followed, follower string, handler followUserHandler) error {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/%s/follow", followed), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("username")
	c.SetParamValues(followed)
	req.Header.Set("Goduit-Client-Username", follower)
	return handler.Follow(c)
}
