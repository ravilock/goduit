package followercentral

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	integrationtests "github.com/ravilock/goduit/integrationTests"
	followerCentralModels "github.com/ravilock/goduit/internal/followerCentral/models"
	followerCentralRepositories "github.com/ravilock/goduit/internal/followerCentral/repositories"
	"github.com/ravilock/goduit/internal/mongo"
	profileManagerRequests "github.com/ravilock/goduit/internal/profileManager/requests"
	profileManagerResponses "github.com/ravilock/goduit/internal/profileManager/responses"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestFollow(t *testing.T) {
	// Connect Mongo DB
	client, err := mongo.ConnectDatabase(viper.GetString("db.url"))
	if err != nil {
		log.Fatalln("Error connecting to database", err)
	}
	followerCentralRepository := followerCentralRepositories.NewFollowerRepository(client)
	serverUrl := viper.GetString("server.url")
	followUserEndpoint := fmt.Sprintf("%s%s", serverUrl, "/api/profiles/")
	httpClient := http.Client{}

	t.Run("Should follow a user", func(t *testing.T) {
		// Arrange
		followedIdentity, _ := integrationtests.MustRegisterUser(t, profileManagerRequests.RegisterPayload{})
		followerIdentity, followerToken := integrationtests.MustRegisterUser(t, profileManagerRequests.RegisterPayload{})
		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s%s", followUserEndpoint, followedIdentity.Username, "/followers"), nil)
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", followerToken))

		// Act
		res, err := httpClient.Do(req)

		// Assert
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)
		followResponse := new(profileManagerResponses.ProfileResponse)
		resBytes, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		err = json.Unmarshal(resBytes, followResponse)
		require.NoError(t, err)
		checkFollowResponse(t, followedIdentity.Username, true, followResponse)
		followerModel, err := followerCentralRepository.IsFollowedBy(context.Background(), followedIdentity.Subject, followerIdentity.Subject)
		require.NoError(t, err)
		checkFollowerModel(t, followedIdentity.Subject, followerIdentity.Subject, followerModel)
	})

	t.Run("Should return HTTP 404 if no user is found", func(t *testing.T) {
		// Arrange
		inexistentUsername := "inexistent-username"
		_, followerToken := integrationtests.MustRegisterUser(t, profileManagerRequests.RegisterPayload{})
		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s%s", followUserEndpoint, inexistentUsername, "/followers"), nil)
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", followerToken))

		// Act
		res, err := httpClient.Do(req)

		// Assert
		require.NoError(t, err)
		require.Equal(t, http.StatusNotFound, res.StatusCode)
	})

	t.Run("Should return HTTP 409 if follower already follows followed", func(t *testing.T) {
		// Arrange
		followedIdentity, _ := integrationtests.MustRegisterUser(t, profileManagerRequests.RegisterPayload{})
		_, followerToken := integrationtests.MustRegisterUser(t, profileManagerRequests.RegisterPayload{})
		integrationtests.MustFollowUser(t, followedIdentity.Username, followerToken)
		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s%s", followUserEndpoint, followedIdentity.Username, "/followers"), nil)
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", followerToken))

		// Act
		res, err := httpClient.Do(req)

		// Assert
		require.NoError(t, err)
		require.Equal(t, http.StatusConflict, res.StatusCode)
	})
}

func checkFollowResponse(t *testing.T, username string, following bool, response *profileManagerResponses.ProfileResponse) {
	t.Helper()
	require.Equal(t, username, response.Profile.Username, "User username should be the same")
	require.Equal(t, following, response.Profile.Following)
	require.Zero(t, response.Profile.Image)
	require.Zero(t, response.Profile.Bio)
}

func checkFollowerModel(t *testing.T, followed, follower string, model *followerCentralModels.Follower) {
	t.Helper()
	require.NotNil(t, model)
	require.Equal(t, followed, *model.Followed, "Wrong followed username")
	require.Equal(t, follower, *model.Follower, "Wrong follower username")
}
