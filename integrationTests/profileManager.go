package integrationtests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"testing"

	"github.com/ravilock/goduit/internal/cookie"
	profileManagerRequests "github.com/ravilock/goduit/internal/profileManager/requests"
	profileManagerResponses "github.com/ravilock/goduit/internal/profileManager/responses"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/internal/identity"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func MustRegisterUser(t *testing.T, registerPayload profileManagerRequests.RegisterPayload) (*identity.Identity, *http.Cookie) {
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
	cookie := CheckCookie(t, res)
	id, err := identity.FromToken(cookie.Value)
	require.NoError(t, err)
	return id, cookie
}

func CheckCookie(t *testing.T, res *http.Response) *http.Cookie {
	cookies := res.Cookies()
	cIndex := slices.IndexFunc(cookies, func(c *http.Cookie) bool {
		return c.Name == cookie.CookieKey
	})
	require.NotEqual(t, -1, cIndex, "Cookie not found")
	cookie := cookies[cIndex]
	return cookie
}
