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
	encryptionkeys "github.com/ravilock/goduit/internal/config/encryptionKeys"
	"github.com/ravilock/goduit/internal/config/mongo"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	if err := encryptionkeys.LoadKeys(); err != nil {
		log.Println("Failed to read encrpytion keys", err)
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
	apiGroup := e.Group("/api")
	apiGroup.GET("/healthcheck", handlers.Healthcheck)
	routers.UsersRouter(apiGroup)
	routers.ProfilesRouter(apiGroup)

	// Start server
	e.Logger.Fatal(e.Start(":9191"))
}
