package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/internal/app"
	"github.com/ravilock/goduit/internal/profileManager/assemblers"
	"github.com/ravilock/goduit/internal/profileManager/models"
	"github.com/ravilock/goduit/internal/profileManager/requests"
)

type authenticator interface {
	Login(ctx context.Context, email, password string) (*models.User, string, error)
	profileUpdater
}

type CookieCreator interface {
	Create(token string) *http.Cookie
}

type loginHandler struct {
	service       authenticator
	cookieService CookieCreator
}

func (h *loginHandler) Login(c echo.Context) error {
	request := new(requests.LoginRequest)
	if err := c.Bind(request); err != nil {
		return api.CouldNotUnmarshalBodyError
	}

	if err := request.Validate(); err != nil {
		return err
	}

	ctx := c.Request().Context()
	user, token, err := h.service.Login(ctx, request.User.Email, request.User.Password)
	if err != nil {
		if appError := new(app.AppError); errors.As(err, &appError) {
			switch appError.ErrorCode {
			case app.UserNotFoundErrorCode:
				fallthrough
			case app.WrongPasswordErrorCode:
				return api.FailedLoginAttempt
			}
		}
		return err
	}

	lastSession := time.Now().UTC().Truncate(time.Millisecond)
	user.LastSession = &lastSession
	if _, err := h.service.UpdateProfile(context.Background(), *user.Email, *user.Username, "", user); err != nil {
		log.Println("Error Updating Last Session", err)
	}

	response := assemblers.UserResponse(user)
	cookie := h.cookieService.Create(token)
	c.SetCookie(cookie)
	return c.JSON(http.StatusOK, response)
}
