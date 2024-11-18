package requests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
	"github.com/stretchr/testify/require"
)

func TestUpdateProfile(t *testing.T) {
	t.Run("Valid request should not return errors", func(t *testing.T) {
		server := mockValidImageURL(t)
		defer server.Close()
		request := generateUpdateProfileRequest(server.URL)
		err := request.Validate()
		require.NoError(t, err)
	})

	t.Run("Email is required", func(t *testing.T) {
		server := mockValidImageURL(t)
		defer server.Close()
		request := generateUpdateProfileRequest(server.URL)
		request.User.Email = ""
		err := request.Validate()
		require.ErrorContains(t, err, api.RequiredFieldError("Email").Error())
	})

	t.Run("Username is required", func(t *testing.T) {
		server := mockValidImageURL(t)
		defer server.Close()
		request := generateUpdateProfileRequest(server.URL)
		request.User.Username = ""
		err := request.Validate()
		require.ErrorContains(t, err, api.RequiredFieldError("Username").Error())
	})

	t.Run("Bio is required", func(t *testing.T) {
		server := mockValidImageURL(t)
		defer server.Close()
		request := generateUpdateProfileRequest(server.URL)
		request.User.Bio = ""
		err := request.Validate()
		require.ErrorContains(t, err, api.RequiredFieldError("Bio").Error())
	})

	t.Run("Image is required", func(t *testing.T) {
		server := mockValidImageURL(t)
		defer server.Close()
		request := generateUpdateProfileRequest(server.URL)
		request.User.Image = ""
		err := request.Validate()
		require.ErrorContains(t, err, api.RequiredFieldError("Image").Error())
	})

	t.Run("Email should not be blank", func(t *testing.T) {
		server := mockValidImageURL(t)
		defer server.Close()
		request := generateUpdateProfileRequest(server.URL)
		request.User.Email = " "
		err := request.Validate()
		require.ErrorContains(t, err, api.RequiredFieldError("Email").Error())
	})

	t.Run("Password should not be blank", func(t *testing.T) {
		server := mockValidImageURL(t)
		defer server.Close()
		request := generateUpdateProfileRequest(server.URL)
		request.User.Password = " "
		err := request.Validate()
		require.ErrorContains(t, err, api.RequiredFieldError("Password").Error())
	})

	t.Run("Username should not be blank", func(t *testing.T) {
		server := mockValidImageURL(t)
		defer server.Close()
		request := generateUpdateProfileRequest(server.URL)
		request.User.Username = " "
		err := request.Validate()
		require.ErrorContains(t, err, api.RequiredFieldError("Username").Error())
	})

	t.Run("Bio should not be blank", func(t *testing.T) {
		server := mockValidImageURL(t)
		defer server.Close()
		request := generateUpdateProfileRequest(server.URL)
		request.User.Bio = " "
		err := request.Validate()
		require.ErrorContains(t, err, api.RequiredFieldError("Bio").Error())
	})

	t.Run("Image should not be blank", func(t *testing.T) {
		server := mockValidImageURL(t)
		defer server.Close()
		request := generateUpdateProfileRequest(server.URL)
		request.User.Image = " "
		err := request.Validate()
		require.ErrorContains(t, err, api.RequiredFieldError("Image").Error())
	})

	t.Run("Password should contain at least 8 chars", func(t *testing.T) {
		server := mockValidImageURL(t)
		defer server.Close()
		request := generateUpdateProfileRequest(server.URL)
		request.User.Password = "pass"
		err := request.Validate()
		require.ErrorContains(t, err, api.InvalidFieldLimit("Password", "min", "8").Error())
	})

	t.Run("Username should contain at least 5 chars", func(t *testing.T) {
		server := mockValidImageURL(t)
		defer server.Close()
		request := generateUpdateProfileRequest(server.URL)
		request.User.Username = "user"
		err := request.Validate()
		require.ErrorContains(t, err, api.InvalidFieldLimit("Username", "min", "5").Error())
	})

	t.Run("Email should contain at most 256 chars", func(t *testing.T) {
		server := mockValidImageURL(t)
		defer server.Close()
		request := generateUpdateProfileRequest(server.URL)
		request.User.Email = randomString(256) + "@hotmail.com"
		err := request.Validate()
		require.ErrorContains(t, err, api.InvalidFieldLimit("Email", "max", "256").Error())
	})

	t.Run("Password should contain at most 72 chars", func(t *testing.T) {
		server := mockValidImageURL(t)
		defer server.Close()
		request := generateUpdateProfileRequest(server.URL)
		request.User.Password = randomString(73)
		err := request.Validate()
		require.ErrorContains(t, err, api.InvalidFieldLimit("Password", "max", "72").Error())
	})

	t.Run("Username should contain at most 255 chars", func(t *testing.T) {
		server := mockValidImageURL(t)
		defer server.Close()
		request := generateUpdateProfileRequest(server.URL)
		request.User.Username = randomString(256)
		err := request.Validate()
		require.ErrorContains(t, err, api.InvalidFieldLimit("Username", "max", "255").Error())
	})

	t.Run("Email should be a valid email", func(t *testing.T) {
		server := mockValidImageURL(t)
		defer server.Close()
		request := generateUpdateProfileRequest(server.URL)
		request.User.Email = "email@hotmail."
		err := request.Validate()
		require.ErrorContains(t, err, api.InvalidFieldError("Email", request.User.Email).Error())
	})

	t.Run("Should not accept image URLs that responds with a myme type header different than image", func(t *testing.T) {
		server := mockInvalidImageURL(t)
		defer server.Close()
		request := generateUpdateProfileRequest(server.URL)
		err := request.Validate()
		require.ErrorContains(t, err, api.InvalidImageURLError(server.URL, echo.MIMEApplicationJSON).Error())
	})
}

func generateUpdateProfileRequest(imageURL string) *UpdateProfileRequest {
	updateProfile := new(UpdateProfileRequest)
	updateProfile.User.Username = "test-username"
	updateProfile.User.Email = "test-email@email.com"
	updateProfile.User.Password = "test-password"
	updateProfile.User.Bio = "Test Bio"
	updateProfile.User.Image = imageURL
	return updateProfile
}

func mockValidImageURL(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(echo.HeaderContentType, "image/png")
		w.WriteHeader(200)
		_, err := w.Write(nil)
		require.NoError(t, err)
	}))
}

func mockInvalidImageURL(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
		w.WriteHeader(200)
		_, err := w.Write(nil)
		require.NoError(t, err)
	}))
}
