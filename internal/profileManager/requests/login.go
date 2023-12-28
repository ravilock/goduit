package requests

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/ravilock/goduit/api/validators"
	"github.com/ravilock/goduit/internal/profileManager/models"
)

type Login struct {
	User struct {
		Email    string `json:"email" validate:"required,notblank,min=5,max=256,email"`
		Password string `json:"password" validate:"required,notblank,min=8,max=72"`
	} `json:"user" validate:"required"`
}

func (r *Login) Model() *models.User {
	return &models.User{
		Email: &r.User.Email,
	}
}

func (r *Login) Validate() error {
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

// TODO: Add validation tests here as well
