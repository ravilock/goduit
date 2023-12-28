package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api/handlers"
)

func UsersRouter(apiGroup *echo.Group) {
	usersGroup := apiGroup.Group("/users")
	usersGroup.POST("", handlers.Register)
	usersGroup.POST("/login", handlers.Login)
}
