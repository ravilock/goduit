package profilemanager

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	integrationtests "github.com/ravilock/goduit/integrationTests"
	followerCentralRepositories "github.com/ravilock/goduit/internal/followerCentral/repositories"
	"github.com/ravilock/goduit/internal/mongo"
	profileManagerRequests "github.com/ravilock/goduit/internal/profileManager/requests"
	profileManagerResponses "github.com/ravilock/goduit/internal/profileManager/responses"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestGetProfile(t *testing.T) {
	// Connect Mongo DB
	client, err := mongo.ConnectDatabase(viper.GetString("db.url"))
	if err != nil {
		log.Fatalln("Error connecting to database", err)
	}
	followerCentralRepository := followerCentralRepositories.NewFollowerRepository(client)
	// profileManagerRepository := profileManagerRepositories.NewUserRepository(client)
	serverUrl := viper.GetString("server.url")
	getProfileEndpoint := fmt.Sprintf("%s%s", serverUrl, "/api/profiles")
	httpClient := http.Client{}

	t.Run("Should get a user profile", func(t *testing.T) {
		// Arrange
		id, _ := integrationtests.MustRegisterUser(t, profileManagerRequests.RegisterPayload{})
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", getProfileEndpoint, id.Username), bytes.NewBuffer([]byte{}))
		require.NoError(t, err)

		// Act
		res, err := httpClient.Do(req)

		// Assert
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)
		getProfileResponse := new(profileManagerResponses.ProfileResponse)
		resBytes, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		err = json.Unmarshal(resBytes, getProfileResponse)
		require.NoError(t, err)
		checkGetProfileResponse(t, id.Username, false, getProfileResponse)
	})

	t.Run("Should return following as true if logged user follows profile", func(t *testing.T) {
		// Arrange
		getProfilfeUserIdentity, _ := integrationtests.MustRegisterUser(t, profileManagerRequests.RegisterPayload{})
		followerIdentity, followerToken := integrationtests.MustRegisterUser(t, profileManagerRequests.RegisterPayload{})
		err := followerCentralRepository.Follow(context.Background(), getProfilfeUserIdentity.Subject, followerIdentity.Subject)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", getProfileEndpoint, getProfilfeUserIdentity.Username), bytes.NewBuffer([]byte{}))
		require.NoError(t, err)
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", followerToken))

		// Act
		res, err := httpClient.Do(req)

		// Assert
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)
		getProfileResponse := new(profileManagerResponses.ProfileResponse)
		resBytes, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		err = json.Unmarshal(resBytes, getProfileResponse)
		require.NoError(t, err)
		checkGetProfileResponse(t, getProfilfeUserIdentity.Username, true, getProfileResponse)
	})

	t.Run("Should return HTTP 404 if no user is found", func(t *testing.T) {
		// Arrange
		inexistentUsername := "inexistent-username"
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", getProfileEndpoint, inexistentUsername), bytes.NewBuffer([]byte{}))
		require.NoError(t, err)

		// Act
		res, err := httpClient.Do(req)

		// Assert
		require.NoError(t, err)
		require.Equal(t, http.StatusNotFound, res.StatusCode)
	})
}

func checkGetProfileResponse(t *testing.T, username string, following bool, response *profileManagerResponses.ProfileResponse) {
	t.Helper()
	require.Equal(t, username, response.Profile.Username, "User username should be the same")
	require.Equal(t, following, response.Profile.Following, "Wrong user following")
	require.Zero(t, response.Profile.Image)
	require.Zero(t, response.Profile.Bio)
}
