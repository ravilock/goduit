package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api/middlewares"
	"github.com/ravilock/goduit/internal/profileManager/handlers"
)

func UsersRouter(apiGroup *echo.Group, profileHandler *handlers.ProfileHandler) {
	usersGroup := apiGroup.Group("/users")
	usersGroup.POST("", profileHandler.Register)
	usersGroup.POST("/login", profileHandler.Login)

	userGroup := apiGroup.Group("/user")
	userGroup.GET("", profileHandler.GetOwnProfile, middlewares.CreateAuthMiddleware(true))
	userGroup.PUT("", profileHandler.UpdateProfile, middlewares.CreateAuthMiddleware(true))
}
