package validators

import (
	"testing"

	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/api/requests"
)

func TestLogin(t *testing.T) {
	t.Run("Email is required", func(t *testing.T) {
		request := assembleValidLoginRequest()
		request.User.Email = ""

		got := Login(request)
		assertError(t, got, api.RequiredFieldError("Email"))
	})
	t.Run("Email must not be blank", func(t *testing.T) {
		request := assembleValidLoginRequest()
		request.User.Email = "   "

		got := Login(request)
		assertError(t, got, api.RequiredFieldError("Email"))
	})
	t.Run("Email should have at most 255 chars", func(t *testing.T) {
		request := assembleValidLoginRequest()
		request.User.Email = randomString(256) + "@test.com"

		got := Login(request)
		assertError(t, got, api.InvalidFieldLength("Email", "max", "255"))
	})
	t.Run("Email should be a valid email", func(t *testing.T) {
		request := assembleValidLoginRequest()
		request.User.Email = "not-valid-email@"

		got := Login(request)
		assertError(t, got, api.InvalidFieldError("Email", request.User.Email))
	})
	t.Run("Password is required", func(t *testing.T) {
		request := assembleValidLoginRequest()
		request.User.Password = ""

		got := Login(request)
		assertError(t, got, api.RequiredFieldError("Password"))
	})
	t.Run("Password must not be blank", func(t *testing.T) {
		request := assembleValidLoginRequest()
		request.User.Password = "   "

		got := Login(request)
		assertError(t, got, api.RequiredFieldError("Password"))
	})
	t.Run("Password should have at least 8 chars", func(t *testing.T) {
		request := assembleValidLoginRequest()
		request.User.Password = "pass"

		got := Login(request)
		assertError(t, got, api.InvalidFieldLength("Password", "min", "8"))
	})
	t.Run("Password should have at most 72 chars", func(t *testing.T) {
		request := assembleValidLoginRequest()
		request.User.Password = randomString(73)

		got := Login(request)
		assertError(t, got, api.InvalidFieldLength("Password", "max", "72"))
	})

}

func assembleValidLoginRequest() *requests.Login {
	request := new(requests.Login)
	request.User.Email = "username@test.test"
	request.User.Password = "username-password"
	return request
}
