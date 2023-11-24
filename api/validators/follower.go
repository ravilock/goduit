package validators

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/ravilock/goduit/api/requests"
)

func Follower(request *requests.Follower) error {
	if err := Validate.Struct(request); err != nil {
		if validationErrors := new(validator.ValidationErrors); errors.As(err, validationErrors) {
			for _, validationError := range *validationErrors {
				return ToHTTP(validationError)
			}
		}
		return err
	}
	return nil
}
