package assemblers

import (
	"github.com/ravilock/goduit/api/requests"
	"github.com/ravilock/goduit/api/responses"
	"github.com/ravilock/goduit/internal/app/dtos"
)

func Register(request *requests.Register) *dtos.User {
	return &dtos.User{
		Username: &request.User.Username,
		Email:    &request.User.Email,
		Password: &request.User.Password,
		Token:    new(string),
		Bio:      new(string),
		Image:    new(string),
	}
}

func Login(request *requests.Login) *dtos.User {
	return &dtos.User{
		Email:    &request.User.Email,
		Password: &request.User.Password,
		Username: new(string),
		Token:    new(string),
		Bio:      new(string),
		Image:    new(string),
	}
}

func GetUser(userEmail *string) *dtos.User {
	return &dtos.User{
		Email:    userEmail,
		Password: new(string),
		Username: new(string),
		Token:    new(string),
		Bio:      new(string),
		Image:    new(string),
	}
}

func UpdateUser(request *requests.UpdateUser) *dtos.User {
	userData := request.User
	dto := &dtos.User{Email: &userData.Email}
	if userData.Password != "" {
		dto.Password = &userData.Password
	}
	if userData.Username != "" {
		dto.Username = &userData.Username
	}
	if userData.Bio != "" {
		dto.Bio = &userData.Bio
	}
	if userData.Image != "" {
		dto.Image = &userData.Image
	}
	return dto
}

func UserResponse(user *dtos.User) *responses.User {
	response := new(responses.User)
	if user.Username != nil {
		response.User.Username = *user.Username
	}
	if user.Email != nil {
		response.User.Email = *user.Email
	}
	if user.Token != nil {
		response.User.Token = *user.Token
	}
	if user.Bio != nil {
		response.User.Bio = *user.Bio
	}
	if user.Image != nil {
		response.User.Image = *user.Image
	}
	return response
}
