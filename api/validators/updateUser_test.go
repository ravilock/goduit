package validators

import (
	"encoding/base64"
	"testing"

	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/api/requests"
)

func TestUpdateUser(t *testing.T) {
	InitValidator()
	t.Run("Username must not be blank", func(t *testing.T) {
		request := assembleValidUpdateUserRequest()
		request.User.Username = "   "

		got := UpdateUser(request)
		assertError(t, got, api.RequiredFieldError("Username"))
	})
	t.Run("Username should have at least 5 chars", func(t *testing.T) {
		request := assembleValidUpdateUserRequest()
		request.User.Username = "user"

		got := UpdateUser(request)
		assertError(t, got, api.InvalidFieldLength("Username", "min", "5"))
	})
	t.Run("Username should have at most 255 chars", func(t *testing.T) {
		request := assembleValidUpdateUserRequest()
		request.User.Username = randomString(256)

		got := UpdateUser(request)
		assertError(t, got, api.InvalidFieldLength("Username", "max", "255"))
	})
	t.Run("Email is required", func(t *testing.T) {
		request := assembleValidUpdateUserRequest()
		request.User.Email = ""

		got := UpdateUser(request)
		assertError(t, got, api.RequiredFieldError("Email"))
	})
	t.Run("Email must not be blank", func(t *testing.T) {
		request := assembleValidUpdateUserRequest()
		request.User.Email = "   "

		got := UpdateUser(request)
		assertError(t, got, api.RequiredFieldError("Email"))
	})
	t.Run("Email should have at most 255 chars", func(t *testing.T) {
		request := assembleValidUpdateUserRequest()
		request.User.Email = randomString(256) + "@test.com"

		got := UpdateUser(request)
		assertError(t, got, api.InvalidFieldLength("Email", "max", "255"))
	})
	t.Run("Email should be a valid email", func(t *testing.T) {
		request := assembleValidUpdateUserRequest()
		request.User.Email = "not-valid-email@"

		got := UpdateUser(request)
		assertError(t, got, api.InvalidFieldError("Email", request.User.Email))
	})
	t.Run("Password must not be blank", func(t *testing.T) {
		request := assembleValidUpdateUserRequest()
		request.User.Password = "   "

		got := UpdateUser(request)
		assertError(t, got, api.RequiredFieldError("Password"))
	})
	t.Run("Password should have at least 8 chars", func(t *testing.T) {
		request := assembleValidUpdateUserRequest()
		request.User.Password = "pass"

		got := UpdateUser(request)
		assertError(t, got, api.InvalidFieldLength("Password", "min", "8"))
	})
	t.Run("Password should have at most 72 chars", func(t *testing.T) {
		request := assembleValidUpdateUserRequest()
		request.User.Password = randomString(73)

		got := UpdateUser(request)
		assertError(t, got, api.InvalidFieldLength("Password", "max", "72"))
	})
	t.Run("Bio must not be blank", func(t *testing.T) {
		request := assembleValidUpdateUserRequest()
		request.User.Bio = "   "

		got := UpdateUser(request)
		assertError(t, got, api.RequiredFieldError("Bio"))
	})
	t.Run("Bio should have at most 255 chars", func(t *testing.T) {
		request := assembleValidUpdateUserRequest()
		request.User.Bio = randomString(256)

		got := UpdateUser(request)
		assertError(t, got, api.InvalidFieldLength("Bio", "max", "255"))
	})
	t.Run("Image must not be blank", func(t *testing.T) {
		request := assembleValidUpdateUserRequest()
		request.User.Image = "   "

		got := UpdateUser(request)
		assertError(t, got, api.RequiredFieldError("Image"))
	})
	t.Run("Image should have at most 65000 chars", func(t *testing.T) {
		request := assembleValidUpdateUserRequest()
		request.User.Image = randomString(65001)

		got := UpdateUser(request)
		assertError(t, got, api.InvalidFieldLength("Image", "max", "65000"))
	})
	t.Run("Image should be valid http url", func(t *testing.T) {
		request := assembleValidUpdateUserRequest()

		got := UpdateUser(request)
		assertNoError(t, got)

		request.User.Image = "https:/update.user.image.com"

		got = UpdateUser(request)
		assertError(t, got, api.InvalidFieldError("Image", "https:/update.user.image.com"))
	})
	t.Run("Image should be valid base64", func(t *testing.T) {
		request := assembleValidUpdateUserRequest()
		request.User.Image = base64.StdEncoding.EncodeToString([]byte(request.User.Image))

		got := UpdateUser(request)
		assertNoError(t, got)

		request.User.Image += "fail"

		got = UpdateUser(request)
		assertError(t, got, api.InvalidFieldError("Image", "aHR0cHM6Ly91cGRhdGUudXNlci5pbWFnZS5jb20=fail"))
	})
}

func assembleValidUpdateUserRequest() *requests.UpdateUser {
	request := new(requests.UpdateUser)
	request.User.Username = "username"
	request.User.Email = "username@test.test"
	request.User.Password = "username-password"
	request.User.Bio = "Bio"
	request.User.Image = "https://update.user.image.com"
	return request
}
