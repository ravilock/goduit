package assemblers

import (
	"github.com/ravilock/goduit/api/requests"
	"github.com/ravilock/goduit/api/responses"
	"github.com/ravilock/goduit/internal/app/dtos"
)

func Register(registerRequest *requests.Register) *dtos.User {
	return &dtos.User{
		Username: registerRequest.User.Username,
		Email:    registerRequest.User.Email,
		Password: registerRequest.User.Password,
		Token:    new(string),
		Bio:      new(string),
		Image:    new(string),
	}
}

func Login(loginRequest *requests.Login) *dtos.User {
	return &dtos.User{
		Email:    loginRequest.User.Email,
		Password: loginRequest.User.Password,
		Username: new(string),
		Token:    new(string),
		Bio:      new(string),
		Image:    new(string),
	}
}

func Response(user *dtos.User) *responses.User {
	response := new(responses.User)
	response.User.Username = user.Username
	response.User.Email = user.Email
	response.User.Token = user.Token
	response.User.Bio = user.Bio
	response.User.Image = user.Image
	return response
}
