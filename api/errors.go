package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

var CouldNotUnmarshalBodyError *echo.HTTPError = echo.NewHTTPError(http.StatusBadRequest, "Could Not Unmarshall Body")

var FailedLoginAttempt *echo.HTTPError = echo.NewHTTPError(http.StatusUnauthorized, "Login failed; Invalid user ID or password.")

func RequiredFieldError(field string) *echo.HTTPError {
	return &echo.HTTPError{
		Code:    http.StatusBadRequest,
		Message: fmt.Sprintf("field '%v' is required", field),
	}
}

func InvalidFieldError(field string, value any) *echo.HTTPError {
	return &echo.HTTPError{
		Code:    http.StatusBadRequest,
		Message: fmt.Sprintf("'%v' is not valid for field '%s'", value, field),
	}
}

func InvalidFieldLength(field string, validationName string, validationSize string) *echo.HTTPError {
	sizeMessage := "short"
	if validationName == "max" {
		sizeMessage = "long"
	}

	return &echo.HTTPError{
		Code:    http.StatusBadRequest,
		Message: fmt.Sprintf("value in field '%s' is too %s. %s=%s", field, sizeMessage, validationName, validationSize),
	}
}
