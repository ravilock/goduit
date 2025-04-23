package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/api/validators"
	"github.com/ravilock/goduit/internal/app"
	"github.com/ravilock/goduit/internal/cookie"
	"github.com/ravilock/goduit/internal/profileManager/models"
	profileManagerRequests "github.com/ravilock/goduit/internal/profileManager/requests"
	profileManagerResponses "github.com/ravilock/goduit/internal/profileManager/responses"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	loginTestUsername = "login-test-username"
	loginTestEmail    = "login.test.email@test.test"
	loginTestPassword = "login-test-password"
)

func TestLogin(t *testing.T) {
	err := validators.InitValidator()
	require.NoError(t, err)
	cookieManager := cookie.NewCookieManager()
	authenticatorMock := newMockAuthenticator(t)
	cookieCreatorMock := NewMockCookieCreator(t)
	handler := LoginHandler{service: authenticatorMock, cookieService: cookieCreatorMock}
	e := echo.New()

	t.Run("Should successfully login", func(t *testing.T) {
		// Arrange
		loginRequest := generateLoginBody()
		expectedUserID := primitive.NewObjectID()
		expectedUsername := "testing-username"
		now := time.Now().UTC().Truncate(time.Millisecond)
		expectedUserModel := &models.User{
			ID:           &expectedUserID,
			Username:     &expectedUsername,
			Email:        &loginRequest.User.Email,
			PasswordHash: &loginRequest.User.Password,
			Bio:          nil,
			Image:        nil,
			CreatedAt:    &now,
			UpdatedAt:    &now,
			LastSession:  &now,
		}
		requestBody, err := json.Marshal(loginRequest)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		expectedToken := "token"
		expectedCookie := cookieManager.Create(expectedToken)
		authenticatorMock.EXPECT().Login(c.Request().Context(), loginRequest.User.Email, loginRequest.User.Password).Return(expectedUserModel, expectedToken, nil).Once()
		authenticatorMock.EXPECT().UpdateProfile(mock.AnythingOfType("context.backgroundCtx"), loginRequest.User.Email, *expectedUserModel.Username, "", mock.AnythingOfType("*models.User")).Return("", nil).Once()
		cookieCreatorMock.EXPECT().Create(expectedToken).Return(expectedCookie)

		// Act
		err = handler.Login(c)

		// Assert
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rec.Code)
		loginResponse := new(profileManagerResponses.User)
		err = json.Unmarshal(rec.Body.Bytes(), loginResponse)
		require.NoError(t, err)
		checkCookie(t, rec, expectedToken)
		checkLoginResponse(t, loginRequest, loginResponse)
	})

	t.Run("Should return 401 if email is not found", func(t *testing.T) {
		// Arrange
		loginRequest := generateLoginBody()
		requestBody, err := json.Marshal(loginRequest)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		authenticatorMock.EXPECT().Login(c.Request().Context(), loginRequest.User.Email, loginRequest.User.Password).Return(nil, "", app.UserNotFoundError(loginRequest.User.Email, nil)).Once()

		// Act
		err = handler.Login(c)

		// Assert
		require.ErrorIs(t, err, api.FailedLoginAttempt)
	})

	t.Run("Should return 401 if wrong password", func(t *testing.T) {
		// Arrange
		loginRequest := generateLoginBody()
		requestBody, err := json.Marshal(loginRequest)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		authenticatorMock.EXPECT().Login(c.Request().Context(), loginRequest.User.Email, loginRequest.User.Password).Return(nil, "", app.WrongPasswordError).Once()

		// Act
		err = handler.Login(c)

		// Assert
		require.ErrorIs(t, err, api.FailedLoginAttempt)
	})
}

func generateLoginBody() *profileManagerRequests.LoginRequest {
	request := new(profileManagerRequests.LoginRequest)
	request.User.Email = loginTestEmail
	request.User.Password = loginTestPassword
	return request
}

func checkLoginResponse(t *testing.T, request *profileManagerRequests.LoginRequest, response *profileManagerResponses.User) {
	t.Helper()
	require.Equal(t, request.User.Email, response.User.Email, "User email should be the same")
	require.Zero(t, response.User.Image)
	require.Zero(t, response.User.Bio)
}

func checkCookie(t *testing.T, rec *httptest.ResponseRecorder, expectedToken string) {
	res, ok := rec.Result().Header["Set-Cookie"]
	require.True(t, ok)
	cookiesString := strings.Join(res, ", ")
	require.Contains(t, cookiesString, expectedToken)
	require.Contains(t, cookiesString, cookie.CookieKey)
}
