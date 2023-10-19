package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/api/requests"
	"github.com/ravilock/goduit/api/responses"
	"github.com/ravilock/goduit/internal/app/models"
	"github.com/ravilock/goduit/internal/app/repositories"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	clearDatabase()
	e := echo.New()
	t.Run("Should create new user", func(t *testing.T) {
		registerRequest := generateRegisterBody()
		requestBody, err := json.Marshal(registerRequest)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = Register(c)
		assert.NoError(t, err)
		if rec.Code != http.StatusCreated {
			t.Errorf("Got status different than %v, got %v", http.StatusCreated, rec.Code)
		}
		registerResponse := new(responses.User)
		err = json.Unmarshal(rec.Body.Bytes(), registerResponse)
		assert.NoError(t, err)
		checkRegisterResponse(t, registerRequest, registerResponse)
		userModel, err := repositories.GetUserByEmail(registerRequest.User.Email, context.Background())
		assert.NoError(t, err)
		checkUserModel(t, registerRequest, userModel)
	})
	t.Run("Should not create user with duplicated email", func(t *testing.T) {
		registerRequest := generateRegisterBody()
		registerRequest.User.Username = "different-username"
		requestBody, err := json.Marshal(registerRequest)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = Register(c)
		assert.ErrorIs(t, err, api.ConfictError)
	})
	t.Run("Should not create user with duplicated username", func(t *testing.T) {
		registerRequest := generateRegisterBody()
		registerRequest.User.Email = "different-email@test.test"
		requestBody, err := json.Marshal(registerRequest)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = Register(c)
		assert.ErrorIs(t, err, api.ConfictError)
	})
}

func generateRegisterBody() *requests.Register {
	request := new(requests.Register)
	request.User.Email = "test.test@test.test"
	request.User.Username = "test-username"
	request.User.Password = "test-password"
	return request
}

func checkRegisterResponse(t *testing.T, request *requests.Register, response *responses.User) {
	t.Helper()
	assert.Equal(t, request.User.Email, response.User.Email, "User email should be the same")
	assert.Equal(t, request.User.Username, response.User.Username, "User Username should be the same")
	assert.NotZero(t, response.User.Token)
	assert.Zero(t, response.User.Image)
	assert.Zero(t, response.User.Bio)
}

func checkUserModel(t *testing.T, request *requests.Register, user *models.User) {
	t.Helper()
	assert.Equal(t, request.User.Email, *user.Email, "User email should be the same")
	assert.Equal(t, request.User.Username, *user.Username, "User Username should be the same")
	assert.NotZero(t, *user.PasswordHash)
	assert.Zero(t, *user.Image)
	assert.Zero(t, *user.Bio)
}
