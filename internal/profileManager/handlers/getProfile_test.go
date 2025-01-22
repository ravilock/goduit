package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/api/validators"
	"github.com/ravilock/goduit/internal/app"
	"github.com/ravilock/goduit/internal/profileManager/models"
	profileManagerResponses "github.com/ravilock/goduit/internal/profileManager/responses"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetProfile(t *testing.T) {
	err := validators.InitValidator()
	require.NoError(t, err)
	profileGetterMock := newMockProfileGetter(t)
	isFollowedCheckerMock := newMockIsFollowedChecker(t)
	handler := getProfileHandler{service: profileGetterMock, followerCentral: isFollowedCheckerMock}
	e := echo.New()

	t.Run("Should get user profile", func(t *testing.T) {
		// Arrange
		expectedUserID := primitive.NewObjectID()
		expectedProfileUsername := "test-username"
		expectedProfileEmail := "test.email@test.test"
		expectedUserPassword := "test-password"
		now := time.Now().UTC().Truncate(time.Millisecond)
		expectedUserModel := &models.User{
			ID:           &expectedUserID,
			Username:     &expectedProfileUsername,
			Email:        &expectedProfileEmail,
			PasswordHash: &expectedUserPassword,
			CreatedAt:    &now,
			UpdatedAt:    &now,
			LastSession:  &now,
		}
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/profiles/%s", expectedProfileUsername), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("username")
		c.SetParamValues(expectedProfileUsername)
		profileGetterMock.EXPECT().GetProfileByUsername(c.Request().Context(), expectedProfileUsername).Return(expectedUserModel, nil).Once()
		isFollowedCheckerMock.EXPECT().IsFollowedBy(c.Request().Context(), expectedUserID.Hex(), "").Return(false).Once()

		// Act
		err := handler.GetProfile(c)

		// Assert
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rec.Code)
		getProfileResponse := new(profileManagerResponses.ProfileResponse)
		err = json.Unmarshal(rec.Body.Bytes(), getProfileResponse)
		require.NoError(t, err)
		checkGetProfileResponse(t, expectedProfileUsername, false, getProfileResponse)
	})

	t.Run("Should return following as true if logged user follows profile", func(t *testing.T) {
		// Arrange
		followerUserID := primitive.NewObjectID()
		followerUsername := "follower-username"
		followerEmail := "follower.email@test.test"
		expectedUserID := primitive.NewObjectID()
		expectedProfileUsername := "test-username"
		expectedProfileEmail := "test.email@test.test"
		expectedUserPassword := "test-password"
		now := time.Now().UTC().Truncate(time.Millisecond)
		expectedUserModel := &models.User{
			ID:           &expectedUserID,
			Username:     &expectedProfileUsername,
			Email:        &expectedProfileEmail,
			PasswordHash: &expectedUserPassword,
			CreatedAt:    &now,
			UpdatedAt:    &now,
			LastSession:  &now,
		}
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/profiles/%s", expectedProfileUsername), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Goduit-Subject", followerUserID.Hex())
		req.Header.Set("Goduit-Client-Username", followerUsername)
		req.Header.Set("Goduit-Client-Email", followerEmail)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("username")
		c.SetParamValues(expectedProfileUsername)
		profileGetterMock.EXPECT().GetProfileByUsername(c.Request().Context(), expectedProfileUsername).Return(expectedUserModel, nil).Once()
		isFollowedCheckerMock.EXPECT().IsFollowedBy(c.Request().Context(), expectedUserID.Hex(), followerUserID.Hex()).Return(true).Once()

		// Act
		err := handler.GetProfile(c)

		// Assert
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rec.Code)
		getProfileResponse := new(profileManagerResponses.ProfileResponse)
		err = json.Unmarshal(rec.Body.Bytes(), getProfileResponse)
		require.NoError(t, err)
		checkGetProfileResponse(t, expectedProfileUsername, true, getProfileResponse)
	})

	t.Run("Should return HTTP 404 if no user is found", func(t *testing.T) {
		// Arrange
		inexistentUsername := "inexistent-username"
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/profiles/%s", inexistentUsername), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("username")
		c.SetParamValues(inexistentUsername)
		profileGetterMock.EXPECT().GetProfileByUsername(c.Request().Context(), inexistentUsername).Return(nil, app.UserNotFoundError(inexistentUsername, nil)).Once()

		// Act
		err := handler.GetProfile(c)

		// Assert
		require.ErrorContains(t, err, api.UserNotFound(inexistentUsername).Error())
	})
}

func checkGetProfileResponse(t *testing.T, username string, following bool, response *profileManagerResponses.ProfileResponse) {
	t.Helper()
	require.Equal(t, username, response.Profile.Username, "User username should be the same")
	require.Equal(t, following, response.Profile.Following, "Wrong user following")
	require.Zero(t, response.Profile.Image)
	require.Zero(t, response.Profile.Bio)
}
