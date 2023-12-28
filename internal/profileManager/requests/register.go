package requests

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/ravilock/goduit/api/validators"
	"github.com/ravilock/goduit/internal/profileManager/models"
)

type Register struct {
	User struct {
		Username string `json:"username" validate:"required,notblank,min=5,max=255"`
		Email    string `json:"email" validate:"required,notblank,max=256,email"`
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

func (r *Register) Validate() error {
	if err := validators.Validate.Struct(r); err != nil {
		if validationErrors := new(validator.ValidationErrors); errors.As(err, validationErrors) {
			for _, validationError := range *validationErrors {
				return validators.ToHTTP(validationError)
			}
		}
		return err
	}
	return nil
}
