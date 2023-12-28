package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

var CouldNotUnmarshalBodyError *echo.HTTPError = echo.NewHTTPError(http.StatusBadRequest, "Could Not Unmarshall Body")

var FailedLoginAttempt *echo.HTTPError = echo.NewHTTPError(http.StatusUnauthorized, "Login failed; Invalid user ID or password.")

var FailedAuthentication *echo.HTTPError = echo.NewHTTPError(http.StatusUnauthorized, "Invalid, Empty or Expired Token")

var ConfictError *echo.HTTPError = echo.NewHTTPError(http.StatusConflict, "Content Already Exists")

func UnexpectedTokenSigningMethod(algName string) *echo.HTTPError {
	return &echo.HTTPError{
		Code:    http.StatusUnauthorized,
		Message: fmt.Sprintf("Unexpected signing method: %v", algName),
	}
}

func RequiredFieldError(field string) *echo.HTTPError {
	return &echo.HTTPError{
		Code:    http.StatusBadRequest,
		Message: fmt.Sprintf("field '%v' is required", field),
	}
}

func RequiredOneOfFields(fields []string) *echo.HTTPError {
	return &echo.HTTPError{
		Code:    http.StatusBadRequest,
		Message: fmt.Sprintf("At least one of %q must be provided", strings.Join(fields, ", ")),
	}
}

func InvalidFieldError(field string, value any) *echo.HTTPError {
	return &echo.HTTPError{
		Code:    http.StatusBadRequest,
		Message: fmt.Sprintf("'%v' is not valid for field '%s'", value, field),
	}
}

func UniqueFieldError(field string) *echo.HTTPError {
	return &echo.HTTPError{
		Code:    http.StatusBadRequest,
		Message: fmt.Sprintf("field '%v' must have unique values", field),
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

func UserNotFound(identifier string) *echo.HTTPError {
	return &echo.HTTPError{
		Code:    http.StatusNotFound,
		Message: fmt.Sprintf("User with identifier %q not found", identifier),
	}
}

func FollowerRelationshipNotFound(followed, follower string) *echo.HTTPError {
	return &echo.HTTPError{
		Code:    http.StatusNotFound,
		Message: fmt.Sprintf("%q does not follow %q", follower, followed),
	}
}

func ArticleNotFound(identifier string) *echo.HTTPError {
	return &echo.HTTPError{
		Code:    http.StatusNotFound,
		Message: fmt.Sprintf("Article with identifier %q not found", identifier),
	}
}

func InternalError(internal error) *echo.HTTPError {
	return &echo.HTTPError{
		Code:     http.StatusInternalServerError,
		Message:  "Internal Server Error",
		Internal: internal,
	}
}
