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

func TestFollow(t *testing.T) {
	validators.InitValidator()
	userFollowerMock := newMockUserFollower(t)
	profileGetterMock := newMockProfileGetter(t)
	handler := followUserHandler{userFollowerMock, profileGetterMock}
	e := echo.New()

	t.Run("Should follow a user", func(t *testing.T) {
		// Arrange
		followerID := primitive.NewObjectID()
		followerUsername := "follower-username"
		followerEmail := "follower.email@test.test"
		followedID := primitive.NewObjectID()
		followedUsername := "followed-test-username"
		followedEmail := "followed.email@test.test"
		now := time.Now()
		followedUser := &models.User{
			ID:          &followedID,
			Username:    &followedUsername,
			Email:       &followedEmail,
			CreatedAt:   &now,
			UpdatedAt:   &now,
			LastSession: &now,
		}
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/%s/follow", followedUsername), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("username")
		c.SetParamValues(followedUsername)
		req.Header.Set("Goduit-Subject", followerID.Hex())
		req.Header.Set("Goduit-Client-Username", followerUsername)
		req.Header.Set("Goduit-Client-Email", followerEmail)
		ctx := c.Request().Context()
		profileGetterMock.EXPECT().GetProfileByUsername(ctx, followedUsername).Return(followedUser, nil).Once()
		userFollowerMock.EXPECT().Follow(ctx, followedID.Hex(), followerID.Hex()).Return(nil).Once()

		// Act
		err := handler.Follow(c)

		// Assert
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rec.Code)
		followResponse := new(profileManagerResponses.ProfileResponse)
		err = json.Unmarshal(rec.Body.Bytes(), followResponse)
		require.NoError(t, err)
		checkFollowResponse(t, followedUsername, true, followResponse)
	})

	t.Run("Should return HTTP 404 if no user is found", func(t *testing.T) {
		// Arrange
		followerID := primitive.NewObjectID()
		followerUsername := "follower-username"
		followerEmail := "follower.email@test.test"
		inexistentUser := "inexistent-username"
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/%s/follow", inexistentUser), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("username")
		c.SetParamValues(inexistentUser)
		req.Header.Set("Goduit-Subject", followerID.Hex())
		req.Header.Set("Goduit-Client-Username", followerUsername)
		req.Header.Set("Goduit-Client-Email", followerEmail)
		ctx := c.Request().Context()
		profileGetterMock.EXPECT().GetProfileByUsername(ctx, inexistentUser).Return(nil, app.UserNotFoundError(inexistentUser, nil)).Once()

		// Act
		err := handler.Follow(c)

		// Assert
		require.ErrorContains(t, err, api.UserNotFound(inexistentUser).Error())
	})

	t.Run("Should return HTTP 409 if follower already follows followed", func(t *testing.T) {
		// Arrange
		followerID := primitive.NewObjectID()
		followerUsername := "follower-username"
		followerEmail := "follower.email@test.test"
		followedID := primitive.NewObjectID()
		followedUsername := "followed-test-username"
		followedEmail := "followed.email@test.test"
		now := time.Now()
		followedUser := &models.User{
			ID:          &followedID,
			Username:    &followedUsername,
			Email:       &followedEmail,
			CreatedAt:   &now,
			UpdatedAt:   &now,
			LastSession: &now,
		}
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/%s/follow", followedUsername), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("username")
		c.SetParamValues(followedUsername)
		req.Header.Set("Goduit-Subject", followerID.Hex())
		req.Header.Set("Goduit-Client-Username", followerUsername)
		req.Header.Set("Goduit-Client-Email", followerEmail)
		ctx := c.Request().Context()
		profileGetterMock.EXPECT().GetProfileByUsername(ctx, followedUsername).Return(followedUser, nil).Once()
		userFollowerMock.EXPECT().Follow(ctx, followedID.Hex(), followerID.Hex()).Return(app.ConflictError("followers")).Once()

		// Act
		err := handler.Follow(c)

		// Assert
		require.ErrorContains(t, err, api.ConfictError.Error())
	})
}

func checkFollowResponse(t *testing.T, username string, following bool, response *profileManagerResponses.ProfileResponse) {
	t.Helper()
	require.Equal(t, username, response.Profile.Username, "User username should be the same")
	require.Equal(t, following, response.Profile.Following)
	require.Zero(t, response.Profile.Image)
	require.Zero(t, response.Profile.Bio)
}
