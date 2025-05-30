package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/internal/app"
	"github.com/ravilock/goduit/internal/profileManager/assemblers"
	"github.com/ravilock/goduit/internal/profileManager/models"
	"github.com/ravilock/goduit/internal/profileManager/requests"
)

type profileRegister interface {
	Register(ctx context.Context, model *models.User, password string) (string, error)
}

type RegisterProfileHandler struct {
	service       profileRegister
	cookieService CookieCreator
}

func NewRegisterProfileHandler(service profileRegister, cookieService CookieCreator) *RegisterProfileHandler {
	return &RegisterProfileHandler{
		service:       service,
		cookieService: cookieService,
	}
}

func (h *RegisterProfileHandler) Register(c echo.Context) error {
	request := new(requests.RegisterRequest)
	if err := c.Bind(request); err != nil {
		return api.CouldNotUnmarshalBodyError
	}

	if err := request.Validate(); err != nil {
		return err
	}

	user := request.Model()

	token, err := h.service.Register(c.Request().Context(), user, request.User.Password)
	if err != nil {
		if appError := new(app.AppError); errors.As(err, &appError) {
			switch appError.ErrorCode {
			case app.ConflictErrorCode:
				return api.ConfictError
			}
		}
		return err
	}

	response := assemblers.UserResponse(user)
	cookie := h.cookieService.Create(token)
	c.SetCookie(cookie)
	return c.JSON(http.StatusCreated, response)
}
