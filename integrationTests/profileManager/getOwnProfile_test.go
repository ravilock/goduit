package profilemanager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	integrationtests "github.com/ravilock/goduit/integrationTests"
	profileManagerRequests "github.com/ravilock/goduit/internal/profileManager/requests"
	profileManagerResponses "github.com/ravilock/goduit/internal/profileManager/responses"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestGetOwnProfile(t *testing.T) {
	serverUrl := viper.GetString("server.url")
	getOwnProfileEndpoint := fmt.Sprintf("%s%s", serverUrl, "/api/user")
	httpClient := http.Client{}
	t.Run("Should get authenticaed user's profile", func(t *testing.T) {
		id, token := integrationtests.MustRegisterUser(t, profileManagerRequests.RegisterPayload{})
		req, err := http.NewRequest(http.MethodGet, getOwnProfileEndpoint, bytes.NewBuffer([]byte{}))
		require.NoError(t, err)
		req.Header.Set(echo.HeaderAuthorization, token)
		res, err := httpClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)
		getOwnProfileResponse := new(profileManagerResponses.User)
		resBytes, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		err = json.Unmarshal(resBytes, getOwnProfileResponse)
		require.NoError(t, err)
		checkGetOwnProfileResponse(t, id.Username, id.UserEmail, getOwnProfileResponse)
	})
	t.Run("Should return http 404 if no user is found", func(t *testing.T) {
	})
}

func checkGetOwnProfileResponse(t *testing.T, username, email string, response *profileManagerResponses.User) {
	t.Helper()
	require.Equal(t, email, response.User.Email, "User email should be the same")
	require.Equal(t, username, response.User.Username, "User username should be the same")
	require.Zero(t, response.User.Image)
	require.Zero(t, response.User.Bio)
}
