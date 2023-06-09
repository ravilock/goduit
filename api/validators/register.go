package validators

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/ravilock/goduit/api/requests"
)

func Register(request *requests.Register) error {
	if err := Validate.Struct(request); err != nil {
		if validationErrors := new(validator.ValidationErrors); errors.As(err, validationErrors) {
			for _, validationError := range *validationErrors {
				return toHTTP(validationError)
			}
		}
		return err
	}
	return nil
}
