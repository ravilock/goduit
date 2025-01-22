package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/internal/profileManager/models"
	profileManagerResponses "github.com/ravilock/goduit/internal/profileManager/responses"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetOwnProfile(t *testing.T) {
	profileGetterMock := newMockProfileGetter(t)
	handler := getOwnProfileHandler{service: profileGetterMock}
	e := echo.New()

	t.Run("Should get user's authenticated profile", func(t *testing.T) {
		// Arrange
		expectedSubject := primitive.NewObjectID()
		clientUsername := "test-username"
		clientEmail := "test.email@test.test"
		expectedUserPassword := "test-password"
		now := time.Now().UTC().Truncate(time.Millisecond)
		expectedUserModel := &models.User{
			ID:           &expectedSubject,
			Username:     &clientUsername,
			Email:        &clientEmail,
			PasswordHash: &expectedUserPassword,
			Bio:          nil,
			Image:        nil,
			CreatedAt:    &now,
			UpdatedAt:    &now,
			LastSession:  &now,
		}
		req := httptest.NewRequest(http.MethodGet, "/user", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", expectedSubject.Hex())
		req.Header.Set("Goduit-Client-Username", clientUsername)
		req.Header.Set("Goduit-Client-Email", clientEmail)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		profileGetterMock.EXPECT().GetProfileByID(c.Request().Context(), expectedSubject.Hex()).Return(expectedUserModel, nil).Once()

		// Act
		err := handler.GetOwnProfile(c)

		// Assert
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rec.Code)
		getOwnProfileResponse := new(profileManagerResponses.User)
		err = json.Unmarshal(rec.Body.Bytes(), getOwnProfileResponse)
		require.NoError(t, err)
		checkGetOwnProfileResponse(t, expectedUserModel, getOwnProfileResponse)
	})
}

func checkGetOwnProfileResponse(t *testing.T, expectedUserData *models.User, response *profileManagerResponses.User) {
	t.Helper()
	require.Equal(t, *expectedUserData.Email, response.User.Email, "User email should be the same")
	require.Equal(t, *expectedUserData.Username, response.User.Username, "User username should be the same")
	require.Zero(t, response.User.Image)
	require.Zero(t, response.User.Bio)
}
