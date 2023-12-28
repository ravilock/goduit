package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api/handlers"
	"github.com/ravilock/goduit/api/middlewares"
)

func UsersRouter(apiGroup *echo.Group) {
	usersGroup := apiGroup.Group("/users")
	usersGroup.POST("", handlers.Register)
	usersGroup.GET("", handlers.GetUser, middlewares.AuthenticationMiddleware)
	usersGroup.POST("/login", handlers.Login)
}
