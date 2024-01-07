package requests

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/api/validators"
	"github.com/ravilock/goduit/internal/profileManager/models"
)

type UpdateProfileRequest struct {
	User UpdateProfileUserPayload `json:"user" validate:"required"`
}

type UpdateProfileUserPayload struct {
	Username string `json:"username" validate:"required,omitempty,notblank,min=5,max=255"`
	Email    string `json:"email" validate:"required,notblank,max=256,email"`
	Password string `json:"password" validate:"omitempty,notblank,min=8,max=72"`
	Bio      string `json:"bio" validate:"required,omitempty,notblank,max=255"`
	Image    string `json:"image" validate:"required,omitempty,notblank,max=65000,http_url|base64"`
}

func (r *UpdateProfileRequest) Model() *models.User {
	model := &models.User{
		Username: &r.User.Username,
		Email:    &r.User.Email,
		Bio:      &r.User.Bio,
		Image:    &r.User.Image,
	}
	return model
}

func (r *UpdateProfileRequest) Validate() error {
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
