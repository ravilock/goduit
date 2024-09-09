package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/api/validators"
	"github.com/ravilock/goduit/internal/app"
	profileManagerRequests "github.com/ravilock/goduit/internal/profileManager/requests"
	profileManagerResponses "github.com/ravilock/goduit/internal/profileManager/responses"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUpdateProfile(t *testing.T) {
	validators.InitValidator()
	profileUpdaterMock := newMockProfileUpdater(t)
	handler := updateProfileHandler{service: profileUpdaterMock}
	e := echo.New()
	imageServer := mockValidImageURL(t)
	defer imageServer.Close()

	t.Run("Should fully update an authenticated user's profile", func(t *testing.T) {
		// Arrange
		expectedSubject := primitive.NewObjectID().Hex()
		clientUsername := "test-username"
		clientEmail := "test.email@test.test"
		updateProfileRequest := generateUpdateProfileBody(imageServer.URL)
		requestBody, err := json.Marshal(updateProfileRequest)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPut, "/user", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", expectedSubject)
		req.Header.Set("Goduit-Client-Username", clientUsername)
		req.Header.Set("Goduit-Client-Email", clientEmail)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		expectedToken := "token"
		profileUpdaterMock.EXPECT().UpdateProfile(c.Request().Context(), clientEmail, clientUsername, updateProfileRequest.User.Password, updateProfileRequest.Model()).Return(expectedToken, nil).Once()

		// Act
		err = handler.UpdateProfile(c)

		// Assert
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rec.Code)
		updateProfileResponse := new(profileManagerResponses.User)
		err = json.Unmarshal(rec.Body.Bytes(), updateProfileResponse)
		require.NoError(t, err)
		checkUpdateProfileResponse(t, updateProfileRequest, updateProfileResponse)
		require.Equal(t, expectedToken, updateProfileResponse.User.Token)
	})

	t.Run("Should return HTTP 409 if new username or email is already being used", func(t *testing.T) {
		// Arrange
		expectedSubject := primitive.NewObjectID().Hex()
		clientUsername := "test-username"
		clientEmail := "test.email@test.test"
		updateProfileRequest := generateUpdateProfileBody(imageServer.URL)
		requestBody, err := json.Marshal(updateProfileRequest)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPut, "/user", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", expectedSubject)
		req.Header.Set("Goduit-Client-Username", clientUsername)
		req.Header.Set("Goduit-Client-Email", clientEmail)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		profileUpdaterMock.EXPECT().UpdateProfile(c.Request().Context(), clientEmail, clientUsername, updateProfileRequest.User.Password, updateProfileRequest.Model()).Return("", app.ConflictError("users")).Once()

		// Act
		err = handler.UpdateProfile(c)

		// Assert
		require.ErrorIs(t, err, api.ConfictError)
	})
}

func generateUpdateProfileBody(imageURL string) *profileManagerRequests.UpdateProfileRequest {
	request := new(profileManagerRequests.UpdateProfileRequest)
	request.User.Username = uuid.NewString()
	request.User.Email = fmt.Sprintf("%s@test.test", request.User.Username)
	request.User.Password = uuid.NewString()
	request.User.Bio = uuid.NewString()
	request.User.Image = imageURL
	return request
}

func checkUpdateProfileResponse(t *testing.T, request *profileManagerRequests.UpdateProfileRequest, response *profileManagerResponses.User) {
	t.Helper()
	require.Equal(t, request.User.Username, response.User.Username, "Updated user's username should be %q, got %q", request.User.Username, response.User.Username)
	require.Equal(t, request.User.Email, response.User.Email, "Updated user's email should be %q, got %q", request.User.Email, response.User.Email)
	require.Equal(t, request.User.Bio, response.User.Bio, "Update user's bio should be %q, got %q", request.User.Bio, response.User.Bio)
	require.Equal(t, request.User.Image, response.User.Image, "Update user's image should be %q, got %q", request.User.Image, response.User.Image)
}

func mockValidImageURL(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(echo.HeaderContentType, "image/png")
		w.WriteHeader(200)
		_, err := w.Write(nil)
		require.NoError(t, err)
	}))
}
