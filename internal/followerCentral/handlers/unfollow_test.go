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
	"github.com/ravilock/goduit/internal/profileManager/models"
	profileManagerResponses "github.com/ravilock/goduit/internal/profileManager/responses"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUnfollow(t *testing.T) {
	userUnfollowerMock := newMockUserUnfollower(t)
	profileGetterMock := newMockProfileGetter(t)
	handler := unfollowUserHandler{userUnfollowerMock, profileGetterMock}
	e := echo.New()
	t.Run("Should unfollow a user", func(t *testing.T) {
		// Act
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
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/%s/unfollow", followedUsername), nil)
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
		userUnfollowerMock.EXPECT().Unfollow(ctx, followerID.Hex(), followedID.Hex()).Return(nil).Once()

		// Act
		err := handler.Unfollow(c)

		// Assert
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rec.Code)
		unfollowResponse := new(profileManagerResponses.ProfileResponse)
		err = json.Unmarshal(rec.Body.Bytes(), unfollowResponse)
		require.NoError(t, err)
		checkFollowResponse(t, followerUsername, false, unfollowResponse)
	})

	t.Run("If the user is already not following the other user, return HTTP 200", func(t *testing.T) {
		// Act
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
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/%s/unfollow", followedUsername), nil)
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
		userUnfollowerMock.EXPECT().Unfollow(ctx, followerID.Hex(), followedID.Hex()).Return(nil).Once()

		// Act
		err := handler.Unfollow(c)

		// Assert
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rec.Code)
		unfollowResponse := new(profileManagerResponses.ProfileResponse)
		err = json.Unmarshal(rec.Body.Bytes(), unfollowResponse)
		require.NoError(t, err)
		checkFollowResponse(t, followerUsername, false, unfollowResponse)
	})

	t.Run("Should return 404 if no user is found", func(t *testing.T) {
		inexistentUsername := "inexistent-username"
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/%s/unfollow", inexistentUsername), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("username")
		c.SetParamValues(inexistentUsername)
		req.Header.Set("Goduit-Subject", followerIdentity.Subject)
		req.Header.Set("Goduit-Client-Username", followerIdentity.Username)
		req.Header.Set("Goduit-Client-Email", followerIdentity.UserEmail)
		err := handler.Unfollow(c)
		require.ErrorContains(t, err, api.UserNotFound(inexistentUsername).Error())
	})
}
