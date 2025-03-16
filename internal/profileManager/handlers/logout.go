package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type CookieClearer interface {
	CookieClear() *http.Cookie
}

type logoutHandler struct {
	service CookieClearer
}

func (h *logoutHandler) Logout(c echo.Context) error {
	cookie := h.service.CookieClear()
	c.SetCookie(cookie)
	return c.NoContent(http.StatusOK)
}
