package requests

import "github.com/ravilock/goduit/internal/app/models"

type Register struct {
	User struct {
		Username string `json:"username" validate:"required,notblank,min=5,max=255"`
		Email    string `json:"email" validate:"required,notblank,max=255,email"`
		Password string `json:"password" validate:"required,notblank,min=8,max=72"`
	} `json:"user" validate:"required"`
}

func (r *Register) Model() *models.User {
	return &models.User{
		Username: &r.User.Username,
		Email:    &r.User.Email,
		Bio:      new(string),
		Image:    new(string),
	}
}
