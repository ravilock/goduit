package integrationtests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	profileManagerResponses "github.com/ravilock/goduit/internal/profileManager/responses"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func MustFollowUser(t *testing.T, followed, followerToken string) {
	httpClient := http.Client{}
	serverUrl := viper.GetString("server.url")
	followUserEndpoint := fmt.Sprintf("%s%s%s%s", serverUrl, "/api/profiles/", followed, "/followers")
	req, err := http.NewRequest(http.MethodPost, followUserEndpoint, nil)
	require.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", followerToken))
	res, err := httpClient.Do(req)
	require.NoError(t, err)
	folllowUserResponse := new(profileManagerResponses.ProfileResponse)
	resBytes, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	err = json.Unmarshal(resBytes, folllowUserResponse)
	require.NoError(t, err)
	require.True(t, folllowUserResponse.Profile.Following)
}
