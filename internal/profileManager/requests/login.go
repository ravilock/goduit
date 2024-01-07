package requests

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/ravilock/goduit/api/validators"
)

type LoginRequest struct {
	User LoginUserPayload `json:"user" validate:"required"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,notblank,max=256,email"`
	Password string `json:"password" validate:"required,notblank,min=8,max=72"`
}

func (r *LoginRequest) Validate() error {
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
