package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/api/validators"
	"github.com/ravilock/goduit/internal/app"
	profileManagerRequests "github.com/ravilock/goduit/internal/profileManager/requests"
	profileManagerResponses "github.com/ravilock/goduit/internal/profileManager/responses"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	err := validators.InitValidator()
	require.NoError(t, err)
	profileRegisterMock := newMockProfileRegister(t)
	handler := registerProfileHandler{service: profileRegisterMock}
	e := echo.New()

	t.Run("Should create new user", func(t *testing.T) {
		// Arrange
		registerRequest := generateRegisterBody()
		requestBody, err := json.Marshal(registerRequest)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		expectedToken := "token"
		profileRegisterMock.EXPECT().Register(c.Request().Context(), registerRequest.Model(), registerRequest.User.Password).Return(expectedToken, nil).Once()

		// Act
		err = handler.Register(c)

		// Assert
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, rec.Code)
		registerResponse := new(profileManagerResponses.User)
		err = json.Unmarshal(rec.Body.Bytes(), registerResponse)
		require.NoError(t, err)
		checkRegisterResponse(t, registerRequest, expectedToken, registerResponse)
	})

	t.Run("Should handle when service returns conflict error", func(t *testing.T) {
		// Arrange
		registerRequest := generateRegisterBody()
		requestBody, err := json.Marshal(registerRequest)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		profileRegisterMock.EXPECT().Register(c.Request().Context(), registerRequest.Model(), registerRequest.User.Password).Return("", app.ConflictError("users")).Once()

		// Act
		err = handler.Register(c)

		// Assert
		require.ErrorIs(t, err, api.ConfictError)
	})
}

func generateRegisterBody() *profileManagerRequests.RegisterRequest {
	request := new(profileManagerRequests.RegisterRequest)
	request.User.Email = "test.test@test.test"
	request.User.Username = "test-username"
	request.User.Password = "test-password"
	return request
}

func checkRegisterResponse(t *testing.T, request *profileManagerRequests.RegisterRequest, expectedToken string, response *profileManagerResponses.User) {
	t.Helper()
	require.Equal(t, request.User.Email, response.User.Email, "User email should be the same")
	require.Equal(t, request.User.Username, response.User.Username, "User Username should be the same")
	require.Equal(t, expectedToken, response.User.Token)
	require.Zero(t, response.User.Image)
	require.Zero(t, response.User.Bio)
}
