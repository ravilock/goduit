package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ravilock/goduit/api/handlers"
	"github.com/ravilock/goduit/api/middlewares"
	"github.com/ravilock/goduit/api/routers"
	"github.com/ravilock/goduit/api/validators"
	"github.com/ravilock/goduit/internal/config/mongo"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	databaseURI := os.Getenv("DB_URL")
	if databaseURI == "" {
		log.Fatal("You must sey your 'DATABASE_URI' environmental variable.")
	}
	// Connect Mongo DB
	mongo.ConnectDatabase(databaseURI)
	defer mongo.DisconnectDatabase()

	// Echo instance
	e := echo.New()

	e.HTTPErrorHandler = middlewares.ErrorMiddleware

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Start Validator
	validators.InitValidator()

	// Routes
	e.GET("/healthcheck", handlers.Healthcheck)
	apiGroup := e.Group("/api")
	routers.UsersRouter(apiGroup)

	// Start server
	e.Logger.Fatal(e.Start(":9191"))
}
