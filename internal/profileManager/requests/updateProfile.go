package requests

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/api/validators"
	"github.com/ravilock/goduit/internal/profileManager/models"
)

type UpdateProfileRequest struct {
	User UpdateProfilePayload `json:"user" validate:"required"`
}

type UpdateProfilePayload struct {
	Username string `json:"username" validate:"required,omitempty,notblank,min=5,max=255"`
	Email    string `json:"email" validate:"required,notblank,max=256,email"`
	Password string `json:"password" validate:"omitempty,notblank,min=8,max=72"`
	Bio      string `json:"bio" validate:"required,notblank,max=255"`
	Image    string `json:"image" validate:"required,notblank,max=65000,http_url"`
}

func (r *UpdateProfileRequest) Model() *models.User {
	updatedAt := time.Now().Truncate(time.Millisecond)
	model := &models.User{
		Username:  &r.User.Username,
		Email:     &r.User.Email,
		Bio:       &r.User.Bio,
		Image:     &r.User.Image,
		UpdatedAt: &updatedAt,
	}
	return model
}

func (r *UpdateProfileRequest) Validate() error {
	if err := validators.Validate.Struct(r); err != nil {
		if validationErrors := new(validator.ValidationErrors); errors.As(err, validationErrors) {
			for _, validationError := range *validationErrors {
				return validators.ToHTTP(validationError)
			}
		}
		return err
	}
	user := r.User
	if user.Username == "" && user.Password == "" && user.Bio == "" && user.Image == "" {
		return api.RequiredOneOfFields([]string{"username", "password", "bio", "image"})
	}
	return checkImageURL(r.User.Image)
}

func checkImageURL(imageURL string) error {
	response, err := http.Get(imageURL)
	if err != nil {
		return err
	}
	contentType := response.Header.Get(echo.HeaderContentType)
	if !strings.Contains(contentType, "image") {
		return api.InvalidImageURLError(imageURL, contentType)
	}
	return nil
}
