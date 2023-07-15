package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api/handlers"
	"github.com/ravilock/goduit/api/middlewares"
)

func UsersRouter(apiGroup *echo.Group) {
	usersGroup := apiGroup.Group("/users")
	usersGroup.POST("", handlers.Register)
	usersGroup.POST("/login", handlers.Login)

	userGroup := apiGroup.Group("/user")
	userGroup.GET("", handlers.GetUser, middlewares.AuthenticationMiddleware)
	userGroup.PUT("", handlers.UpdateUser, middlewares.AuthenticationMiddleware)
}
