package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// healthcheck
func Healthcheck(c echo.Context) error {
	return c.String(http.StatusOK, fmt.Sprintln("OK"))
}
