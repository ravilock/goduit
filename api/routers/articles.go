package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/ravilock/goduit/api/handlers"
	"github.com/ravilock/goduit/api/middlewares"
)

func ArticlesRouter(apiGroup *echo.Group) {
	articlesGroup := apiGroup.Group("/articles")
	articlesGroup.POST("", handlers.CreateArticle, middlewares.CreateAuthMiddleware(true))
}
