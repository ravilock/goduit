package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
)

var Validate *validator.Validate

func InitValidator() error {
	Validate = validator.New()
	return Validate.RegisterValidation("notblank", validators.NotBlank)
}

func ToHTTP(err validator.FieldError) *echo.HTTPError {
	tag := err.Tag()
	switch tag {
	case "required", "notblank":
		return api.RequiredFieldError(err.Field())
	case "min", "max":
		return api.InvalidFieldLength(err.Field(), tag, err.Param())
	case "email", "http_url|base64":
		return api.InvalidFieldError(err.Field(), err.Value())
	case "unique":
		return api.UniqueFieldError(err.Field())
	default:
		return nil
	}
}
