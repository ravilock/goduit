package app

import (
	"fmt"
)

type ErrorCode uint

const (
	UserNotFoundErrorCode ErrorCode = iota + 1
	ArticleNotFoundErrorCode
	WrongPasswordErrorCode
	ConflictErrorCode
)

type AppError struct {
	ErrorCode     ErrorCode
	CustomMessage string
	OriginalError error
}

func (err *AppError) Error() string {
	return fmt.Sprintf("ErrorCode: %d: ErrorMessage: %s: %s", err.ErrorCode, err.CustomMessage, err.OriginalError)
}

func (err *AppError) AddContext(originalError error) *AppError {
	err.OriginalError = originalError
	return err
}

func UserNotFoundError(identifier string, originalError error) *AppError {
	return &AppError{
		ErrorCode:     UserNotFoundErrorCode,
		CustomMessage: fmt.Sprintf("User with identifier %q was not found", identifier),
		OriginalError: originalError,
	}
}

func ArticleNotFoundError(identifier string, originalError error) *AppError {
	return &AppError{
		ErrorCode:     ArticleNotFoundErrorCode,
		CustomMessage: fmt.Sprintf("Article with identifier %q was not found", identifier),
		OriginalError: originalError,
	}
}

func ConflictError(resource string) *AppError {
	return &AppError{
		ErrorCode:     ConflictErrorCode,
		CustomMessage: fmt.Sprintf("Conflict on resource %s", resource),
		OriginalError: nil,
	}
}

var WrongPasswordError *AppError = &AppError{
	ErrorCode:     WrongPasswordErrorCode,
	CustomMessage: "Comparison between password and stored password hash failed",
	OriginalError: nil,
}
