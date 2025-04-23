package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type CookieClearer interface {
	CookieClear() *http.Cookie
}

type LogoutHandler struct {
	service CookieClearer
}

func NewLogoutHandler(service CookieClearer) *LogoutHandler {
	return &LogoutHandler{
		service: service,
	}
}

func (h *LogoutHandler) Logout(c echo.Context) error {
	cookie := h.service.CookieClear()
	c.SetCookie(cookie)
	return c.NoContent(http.StatusOK)
}
