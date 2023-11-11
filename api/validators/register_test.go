package validators

import (
	"testing"

	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/api/requests"
)

func TestRegister(t *testing.T) {
	t.Run("Username is required", func(t *testing.T) {
		request := assembleValidRegisterRequest()
		request.User.Username = ""

		got := Register(request)
		assertError(t, got, api.RequiredFieldError("Username"))
	})
	t.Run("Username must not be blank", func(t *testing.T) {
		request := assembleValidRegisterRequest()
		request.User.Username = "   "

		got := Register(request)
		assertError(t, got, api.RequiredFieldError("Username"))
	})
	t.Run("Username should have at least 5 chars", func(t *testing.T) {
		request := assembleValidRegisterRequest()
		request.User.Username = "user"

		got := Register(request)
		assertError(t, got, api.InvalidFieldLength("Username", "min", "5"))
	})
	t.Run("Username should have at most 255 chars", func(t *testing.T) {
		request := assembleValidRegisterRequest()
		request.User.Username = randomString(256)

		got := Register(request)
		assertError(t, got, api.InvalidFieldLength("Username", "max", "255"))
	})
	t.Run("Email is required", func(t *testing.T) {
		request := assembleValidRegisterRequest()
		request.User.Email = ""

		got := Register(request)
		assertError(t, got, api.RequiredFieldError("Email"))
	})
	t.Run("Email must not be blank", func(t *testing.T) {
		request := assembleValidRegisterRequest()
		request.User.Email = "   "

		got := Register(request)
		assertError(t, got, api.RequiredFieldError("Email"))
	})
	t.Run("Email should have at most 255 chars", func(t *testing.T) {
		request := assembleValidRegisterRequest()
		request.User.Email = randomString(256) + "@test.com"

		got := Register(request)
		assertError(t, got, api.InvalidFieldLength("Email", "max", "255"))
	})
	t.Run("Email should be a valid email", func(t *testing.T) {
		request := assembleValidRegisterRequest()
		request.User.Email = "not-valid-email@"

		got := Register(request)
		assertError(t, got, api.InvalidFieldError("Email", request.User.Email))
	})
	t.Run("Password is required", func(t *testing.T) {
		request := assembleValidRegisterRequest()
		request.User.Password = ""

		got := Register(request)
		assertError(t, got, api.RequiredFieldError("Password"))
	})
	t.Run("Password must not be blank", func(t *testing.T) {
		request := assembleValidRegisterRequest()
		request.User.Password = "   "

		got := Register(request)
		assertError(t, got, api.RequiredFieldError("Password"))
	})
	t.Run("Password should have at least 8 chars", func(t *testing.T) {
		request := assembleValidRegisterRequest()
		request.User.Password = "pass"

		got := Register(request)
		assertError(t, got, api.InvalidFieldLength("Password", "min", "8"))
	})
	t.Run("Password should have at most 72 chars", func(t *testing.T) {
		request := assembleValidRegisterRequest()
		request.User.Password = randomString(73)

		got := Register(request)
		assertError(t, got, api.InvalidFieldLength("Password", "max", "72"))
	})

}

func assembleValidRegisterRequest() *requests.Register {
	request := new(requests.Register)
	request.User.Username = "username"
	request.User.Email = "username@test.test"
	request.User.Password = "username-password"
	return request
}
