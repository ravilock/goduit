package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api/requests"
	"github.com/ravilock/goduit/api/responses"
	"github.com/stretchr/testify/assert"
)

const oldUpdateUserTestUsername = "update-user-test-username"
const updateUserTestUsername = "update-user-test-username"
const updateUserTestEmail = "update.user.email@test.test"
const updateUserTestPassword = "update-user-test-password"
const updateUserTestBio = "update user test bio"
const updateUserTestImage = "https://update.user.test.bio.com/image"

func TestUpdateUser(t *testing.T) {
	clearDatabase()
	e := echo.New()
	if err := registerAccount(oldUpdateUserTestUsername, updateUserTestEmail, ""); err != nil {
		log.Fatal("Could not create user", err)
	}
	t.Run("Should fully update currently authenticated in user", func(t *testing.T) {
		updateUserRequest := generateUpdateUserBody()
		requestBody, err := json.Marshal(updateUserRequest)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPut, "/user", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", updateUserTestEmail)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = UpdateUser(c)
		assert.NoError(t, err)
		if rec.Code != http.StatusOK {
			t.Errorf("Got status different than %v, got %v", http.StatusOK, rec.Code)
		}
		updateUserResponse := new(responses.User)
		err = json.Unmarshal(rec.Body.Bytes(), updateUserResponse)
		assert.NoError(t, err)
		checkUpdateUserResponse(t, updateUserRequest, updateUserResponse)
	})
	t.Run("Should update only the requested fields", func(t *testing.T) {
		request := new(requests.UpdateUser)
		request.User.Username = oldUpdateUserTestUsername
		request.User.Email = updateUserTestEmail
		requestBody, err := json.Marshal(request)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPut, "/user", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", updateUserTestEmail)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = UpdateUser(c)
		assert.NoError(t, err)
		if rec.Code != http.StatusOK {
			t.Errorf("Got status different than %v, got %v", http.StatusOK, rec.Code)
		}
		updateUserResponse := new(responses.User)
		err = json.Unmarshal(rec.Body.Bytes(), updateUserResponse)
		assert.NoError(t, err)
		assert.Equal(t, request.User.Username, updateUserResponse.User.Username, "User email should be the same")
		assert.Equal(t, updateUserTestEmail, updateUserResponse.User.Email, "User email should be the same")
		assert.Equal(t, "", updateUserResponse.User.Bio, "User username should be the same")
		assert.Equal(t, "", updateUserResponse.User.Image, "User username should be the same")
	})
	// TODO: Add test for changing user password
	// TODO: If user updates email successfully, should re-generate token with new email
	// TODO: Check if possible to instantiate echo instance with middlewares (like auth middleware) and do full testing
}

func generateUpdateUserBody() *requests.UpdateUser {
	request := new(requests.UpdateUser)
	request.User.Username = updateUserTestUsername
	request.User.Email = updateUserTestEmail
	request.User.Password = updateUserTestPassword
	request.User.Bio = updateUserTestBio
	request.User.Image = updateUserTestImage
	return request
}

func checkUpdateUserResponse(t *testing.T, request *requests.UpdateUser, response *responses.User) {
	t.Helper()
	assert.Equal(t, request.User.Email, response.User.Email, "User email should be the same")
	assert.Equal(t, request.User.Username, response.User.Username, "User username should be the same")
	assert.Equal(t, request.User.Bio, response.User.Bio, "User username should be the same")
	assert.Equal(t, request.User.Image, response.User.Image, "User username should be the same")
}
