package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/api/requests"
	"github.com/ravilock/goduit/api/responses"
	"github.com/stretchr/testify/assert"
)

const loginTestUsername = "login-test-username"
const loginTestEmail = "login.test.email@test.test"
const loginTestPassword = "login-test-password"

func TestLogin(t *testing.T) {
	clearDatabase()
	if err := registerAccount(loginTestUsername, loginTestEmail, loginTestPassword); err != nil {
		log.Fatal("Could not create user", err)
	}
	e := echo.New()
	t.Run("Should successfully login", func(t *testing.T) {
		loginRequest := generateLoginBody()
		requestBody, err := json.Marshal(loginRequest)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = Login(c)
		assert.NoError(t, err)
		if rec.Code != http.StatusOK {
			t.Errorf("Got status different than %v, got %v", http.StatusOK, rec.Code)
		}
		loginResponse := new(responses.User)
		err = json.Unmarshal(rec.Body.Bytes(), loginResponse)
		assert.NoError(t, err)
		checkLoginResponse(t, loginRequest, loginResponse)
	})
	t.Run("Should return 401 if email is not found", func(t *testing.T) {
		loginRequest := generateLoginBody()
		loginRequest.User.Email = "wrong-email@test.test"
		requestBody, err := json.Marshal(loginRequest)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = Login(c)
		assert.ErrorIs(t, err, api.FailedLoginAttempt)
	})
	t.Run("Should return 401 if wrong password", func(t *testing.T) {
		loginRequest := generateLoginBody()
		loginRequest.User.Password = "wrong-user-password"
		requestBody, err := json.Marshal(loginRequest)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = Login(c)
		assert.ErrorIs(t, err, api.FailedLoginAttempt)
	})
}

func generateLoginBody() *requests.Login {
	request := new(requests.Login)
	request.User.Email = loginTestEmail
	request.User.Password = loginTestPassword
	return request
}

func checkLoginResponse(t *testing.T, request *requests.Login, response *responses.User) {
	t.Helper()
	assert.Equal(t, request.User.Email, response.User.Email, "User email should be the same")
	assert.NotZero(t, response.User.Token)
	assert.Zero(t, response.User.Image)
	assert.Zero(t, response.User.Bio)
}

func login(email, password string) (string, error) {
	if email == "" {
		email = "default.email@test.test"
	}
	if password == "" {
		password = "default-password"
	}
	loginRequest := new(requests.Login)
	loginRequest.User.Email = email
	loginRequest.User.Password = password
	requestBody, err := json.Marshal(loginRequest)
	if err != nil {
		return "", err
	}
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if err := Login(c); err != nil {
		return "", err
	}
	loginResponse := new(responses.User)
	if err = json.Unmarshal(rec.Body.Bytes(), loginResponse); err != nil {
		return "", err
	}
	return loginResponse.User.Token, nil

}