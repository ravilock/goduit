package requests

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/api/validators"
	"github.com/ravilock/goduit/internal/profileManager/models"
)

type UpdateProfile struct {
	User struct {
		Username string `json:"username" validate:"omitempty,notblank,min=5,max=255"`
		Email    string `json:"email" validate:"required,notblank,min=5,max=255,email"`
		Password string `json:"password" validate:"omitempty,notblank,min=8,max=72"`
		Bio      string `json:"bio" validate:"omitempty,notblank,max=255"`
		Image    string `json:"image" validate:"omitempty,notblank,max=65000,http_url|base64"`
	} `json:"user" validate:"required"`
}

func (r *UpdateProfile) Model() *models.User {
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

func (r *UpdateProfile) Validate() error {
	if err := validators.Validate.Struct(r); err != nil {
		if validationErrors := new(validator.ValidationErrors); errors.As(err, validationErrors) {
			for _, validationError := range *validationErrors {
				return validators.ToHTTP(validationError)
			}
		}
		return err
	}
	user := r.User
	if user.Username == "" && user.Password == "" && user.Bio == "" && user.Image == "" {
		return api.RequiredOneOfFields([]string{"username", "password", "bio", "image"})
	}
	return nil
}
