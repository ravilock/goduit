package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
)

var Validate *validator.Validate

func InitValidator() {
	Validate = validator.New()
	Validate.RegisterValidation("notblank", validators.NotBlank)
}

func toHTTP(err validator.FieldError) *echo.HTTPError {
	tag := err.Tag()
	switch tag {
	case "required", "notblank":
		return api.RequiredFieldError(err.Field())
	case "min", "max":
		return api.InvalidFieldLength(err.Field(), tag, err.Param())
	default:
		return nil
	}
}
