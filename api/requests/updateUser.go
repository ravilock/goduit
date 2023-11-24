package requests

import "github.com/ravilock/goduit/internal/app/models"

type UpdateUser struct {
	User struct {
		Username string `json:"username" validate:"omitempty,notblank,min=5,max=255"`
		Email    string `json:"email" validate:"required,notblank,min=5,max=255,email"`
		Password string `json:"password" validate:"omitempty,notblank,min=8,max=72"`
		Bio      string `json:"bio" validate:"omitempty,notblank,max=255"`
		Image    string `json:"image" validate:"omitempty,notblank,max=65000,http_url|base64"`
	} `json:"user" validate:"required"`
}

func (r *UpdateUser) Model() *models.User {
	userData := r.User
	model := &models.User{Email: &userData.Email}
	if userData.Username != "" {
		model.Username = &userData.Username
	}
	if userData.Bio != "" {
		model.Bio = &userData.Bio
	}
	if userData.Image != "" {
		model.Image = &userData.Image
	}
	return model
}
