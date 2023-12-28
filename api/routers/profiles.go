package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api/handlers"
	"github.com/ravilock/goduit/api/middlewares"
)

func ProfilesRouter(apiGroup *echo.Group) {
	profileGroup := apiGroup.Group("/profile")
	profileGroup.GET("/:username", handlers.GetProfile, middlewares.CreateAuthMiddleware(false))
	profileGroup.POST("/:username/follow", handlers.Follow, middlewares.CreateAuthMiddleware(false))
	profileGroup.POST("/:username/unfollow", handlers.Unfollow, middlewares.CreateAuthMiddleware(false))
}
