package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api/responses"
	"github.com/ravilock/goduit/internal/config/mongo"
	"github.com/ravilock/goduit/internal/profileManager/repositories"
	"github.com/ravilock/goduit/internal/profileManager/services"
	"github.com/stretchr/testify/assert"
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
	repository := repositories.NewUserRepository(client)
	manager := services.NewProfileManager(repository)
	handler := NewProfileHandler(manager)
	clearDatabase(client)
	e := echo.New()
	if err := registerUser(getOwnProfileTestUsername, getOwnProfileTestEmail, "", handler.registerProfileHandler); err != nil {
		log.Fatal("Could not create user", err)
	}
	t.Run("Should get user's authenticated profile", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", getOwnProfileTestEmail)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err := handler.GetOwnProfile(c)
		assert.NoError(t, err)
		if rec.Code != http.StatusOK {
			t.Errorf("Got status different than %v, got %v", http.StatusOK, rec.Code)
		}
		getOwnProfileResponse := new(responses.User)
		err = json.Unmarshal(rec.Body.Bytes(), getOwnProfileResponse)
		assert.NoError(t, err)
		checkGetOwnProfileResponse(t, getOwnProfileTestUsername, getOwnProfileTestEmail, getOwnProfileResponse)
	})
}

func checkGetOwnProfileResponse(t *testing.T, username, email string, response *responses.User) {
	t.Helper()
	assert.Equal(t, email, response.User.Email, "User email should be the same")
	assert.Equal(t, username, response.User.Username, "User username should be the same")
	assert.Zero(t, response.User.Image)
	assert.Zero(t, response.User.Bio)
}
