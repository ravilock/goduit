package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/internal/config/mongo"
	followerCentralRepositories "github.com/ravilock/goduit/internal/followerCentral/repositories"
	followerCentral "github.com/ravilock/goduit/internal/followerCentral/services"
	profileManagerRepositories "github.com/ravilock/goduit/internal/profileManager/repositories"
	profileManagerResponses "github.com/ravilock/goduit/internal/profileManager/responses"
	profileManager "github.com/ravilock/goduit/internal/profileManager/services"
	"github.com/stretchr/testify/require"
)

const getOwnProfileTestUsername = "get-own-profile-test-username"
const getOwnProfileTestEmail = "get.own.profile.email@test.test"

func TestGetOwnProfile(t *testing.T) {
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
	t.Run("Should get user's authenticated profile", func(t *testing.T) {
		identity, err := registerUser(getOwnProfileTestUsername, getOwnProfileTestEmail, "", handler.registerProfileHandler)
		require.NoError(t, err, "Could Not Create User")
		req := httptest.NewRequest(http.MethodGet, "/user", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", identity.Subject)
		req.Header.Set("Goduit-Client-Username", identity.Username)
		req.Header.Set("Goduit-Client-Email", identity.UserEmail)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = handler.GetOwnProfile(c)
		require.NoError(t, err)
		if rec.Code != http.StatusOK {
			t.Errorf("Got status different than %v, got %v", http.StatusOK, rec.Code)
		}
		getOwnProfileResponse := new(profileManagerResponses.User)
		err = json.Unmarshal(rec.Body.Bytes(), getOwnProfileResponse)
		require.NoError(t, err)
		checkGetOwnProfileResponse(t, getOwnProfileTestUsername, getOwnProfileTestEmail, getOwnProfileResponse)
	})
}

func checkGetOwnProfileResponse(t *testing.T, username, email string, response *profileManagerResponses.User) {
	t.Helper()
	require.Equal(t, email, response.User.Email, "User email should be the same")
	require.Equal(t, username, response.User.Username, "User username should be the same")
	require.Zero(t, response.User.Image)
	require.Zero(t, response.User.Bio)
}
