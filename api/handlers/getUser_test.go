package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api/responses"
	"github.com/stretchr/testify/assert"
)

const getUserTestUsername = "get-user-test-username"
const getUserTestEmail = "get.user.email@test.test"

func TestGetUser(t *testing.T) {
	clearDatabase()
	e := echo.New()
	if err := registerAccount(getUserTestUsername, getUserTestEmail, ""); err != nil {
		log.Fatal("Could not create user", err)
	}
	t.Run("Should get currently authenticated in user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", getUserTestEmail)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err := GetUser(c)
		assert.NoError(t, err)
		if rec.Code != http.StatusOK {
			t.Errorf("Got status different than %v, got %v", http.StatusOK, rec.Code)
		}
		getUserResponse := new(responses.User)
		err = json.Unmarshal(rec.Body.Bytes(), getUserResponse)
		assert.NoError(t, err)
		checkGetUserResponse(t, getUserTestUsername, getUserTestEmail, getUserResponse)
	})
	// TODO: Check if possible to instantiate echo instance with middlewares (like auth middleware) and do full testing
}

func checkGetUserResponse(t *testing.T, username, email string, response *responses.User) {
	t.Helper()
	assert.Equal(t, email, response.User.Email, "User email should be the same")
	assert.Equal(t, username, response.User.Username, "User username should be the same")
	assert.Zero(t, response.User.Image)
	assert.Zero(t, response.User.Bio)
}
