package validators

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/api/requests"
)

func UpdateUser(request *requests.UpdateUser) error {
	if err := Validate.Struct(request); err != nil {
		if validationErrors := new(validator.ValidationErrors); errors.As(err, validationErrors) {
			for _, validationError := range *validationErrors {
				return toHTTP(validationError)
			}
		}
		return err
	}
	user := request.User
	if user.Username == "" && user.Password == "" && user.Bio == "" && user.Image == "" {
		return api.RequiredOneOfFields([]string{"username", "password", "bio", "image"})
	}
	return nil
}
