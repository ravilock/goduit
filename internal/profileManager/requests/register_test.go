package requests

import (
	"testing"

	"github.com/ravilock/goduit/api"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	t.Run("Valid request should not return errors", func(t *testing.T) {
		request := generateRegisterRequest()
		err := request.Validate()
		assert.NoError(t, err)
	})

	t.Run("Email is required", func(t *testing.T) {
		request := generateRegisterRequest()
		request.User.Email = ""
		err := request.Validate()
		assert.ErrorContains(t, err, api.RequiredFieldError("Email").Error())
	})

	t.Run("Password is required", func(t *testing.T) {
		request := generateRegisterRequest()
		request.User.Password = ""
		err := request.Validate()
		assert.ErrorContains(t, err, api.RequiredFieldError("Password").Error())
	})

	t.Run("Username is required", func(t *testing.T) {
		request := generateRegisterRequest()
		request.User.Username = ""
		err := request.Validate()
		assert.ErrorContains(t, err, api.RequiredFieldError("Username").Error())
	})

	t.Run("Email should not be blank", func(t *testing.T) {
		request := generateRegisterRequest()
		request.User.Email = " "
		err := request.Validate()
		assert.ErrorContains(t, err, api.RequiredFieldError("Email").Error())
	})

	t.Run("Password should not be blank", func(t *testing.T) {
		request := generateRegisterRequest()
		request.User.Password = " "
		err := request.Validate()
		assert.ErrorContains(t, err, api.RequiredFieldError("Password").Error())
	})

	t.Run("Username should not be blank", func(t *testing.T) {
		request := generateRegisterRequest()
		request.User.Username = " "
		err := request.Validate()
		assert.ErrorContains(t, err, api.RequiredFieldError("Username").Error())
	})

	t.Run("Password should contain at least 8 chars", func(t *testing.T) {
		request := generateRegisterRequest()
		request.User.Password = "pass"
		err := request.Validate()
		assert.ErrorContains(t, err, api.InvalidFieldLength("Password", "min", "8").Error())
	})

	t.Run("Username should contain at least 5 chars", func(t *testing.T) {
		request := generateRegisterRequest()
		request.User.Username = "user"
		err := request.Validate()
		assert.ErrorContains(t, err, api.InvalidFieldLength("Username", "min", "5").Error())
	})

	t.Run("Email should contain at most 256 chars", func(t *testing.T) {
		request := generateRegisterRequest()
		request.User.Email = randomString(256) + "@hotmail.com"
		err := request.Validate()
		assert.ErrorContains(t, err, api.InvalidFieldLength("Email", "max", "256").Error())
	})

	t.Run("Password should contain at most 72 chars", func(t *testing.T) {
		request := generateRegisterRequest()
		request.User.Password = randomString(73)
		err := request.Validate()
		assert.ErrorContains(t, err, api.InvalidFieldLength("Password", "max", "72").Error())
	})

	t.Run("Username should contain at most 255 chars", func(t *testing.T) {
		request := generateRegisterRequest()
		request.User.Username = randomString(256)
		err := request.Validate()
		assert.ErrorContains(t, err, api.InvalidFieldLength("Username", "max", "255").Error())
	})

	t.Run("Email should be a valid email", func(t *testing.T) {
		request := generateRegisterRequest()
		request.User.Email = "email@hotmail."
		err := request.Validate()
		assert.ErrorContains(t, err, api.InvalidFieldError("Email", request.User.Email).Error())
	})
}

func generateRegisterRequest() *Register {
	register := new(Register)
	register.User.Email = "test.email@test.com"
	register.User.Username = "test.username"
	register.User.Password = "password123456"
	return register
}