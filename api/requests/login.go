package requests

import "github.com/ravilock/goduit/internal/app/models"

type Login struct {
	User struct {
		Email    string `json:"email" validate:"required,notblank,min=5,max=255,email"`
		Password string `json:"password" validate:"required,notblank,min=8,max=72"`
	} `json:"user" validate:"required"`
}

func (r *Login) Model() *models.User {
	return &models.User{
		Email: &r.User.Email,
	}
}
