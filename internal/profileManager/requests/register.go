package requests

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/ravilock/goduit/api/validators"
	"github.com/ravilock/goduit/internal/profileManager/models"
)

type RegisterRequest struct {
	User RegisterPayload `json:"user" validate:"required"`
}

type RegisterPayload struct {
	Username string `json:"username" validate:"required,notblank,min=5,max=255"`
	Email    string `json:"email" validate:"required,notblank,max=256,email"`
	Password string `json:"password" validate:"required,notblank,min=8,max=72"`
}

func (r *RegisterRequest) Model() *models.User {
	createdAt := time.Now().Truncate(time.Millisecond)
	return &models.User{
		Username:    &r.User.Username,
		Email:       &r.User.Email,
		Bio:         new(string),
		Image:       new(string),
		CreatedAt:   &createdAt,
		UpdatedAt:   nil,
		LastSession: &createdAt,
	}
}

func (r *RegisterRequest) Validate() error {
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
