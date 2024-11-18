package integrationtests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	profileManagerRequests "github.com/ravilock/goduit/internal/profileManager/requests"
	profileManagerResponses "github.com/ravilock/goduit/internal/profileManager/responses"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/internal/identity"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func MustRegisterUser(t *testing.T, registerPayload profileManagerRequests.RegisterPayload) (*identity.Identity, string) {
	httpClient := http.Client{}
	serverUrl := viper.GetString("server.url")
	registerEndpoint := fmt.Sprintf("%s%s", serverUrl, "/api/users")
	if registerPayload.Username == "" {
		registerPayload.Username = UniqueUsername()
	}
	if registerPayload.Email == "" {
		registerPayload.Email = UniqueEmail()
	}
	if registerPayload.Password == "" {
		registerPayload.Password = "12345678"
	}
	requestBody, err := json.Marshal(&profileManagerRequests.RegisterRequest{
		User: registerPayload,
	})
	require.NoError(t, err)
	req, err := http.NewRequest(http.MethodPost, registerEndpoint, bytes.NewBuffer(requestBody))
	require.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res, err := httpClient.Do(req)
	require.NoError(t, err)
	defer res.Body.Close()
	registerResponse := new(profileManagerResponses.User)
	resBytes, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.StatusCode)
	err = json.Unmarshal(resBytes, registerResponse)
	require.NoError(t, err)
	id, err := identity.FromToken(registerResponse.User.Token)
	require.NoError(t, err)
	return id, registerResponse.User.Token
}
