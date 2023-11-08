package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ravilock/goduit/api/handlers"
	"github.com/ravilock/goduit/api/routers"
	"github.com/ravilock/goduit/api/validators"
	encryptionkeys "github.com/ravilock/goduit/internal/config/encryptionKeys"
	"github.com/ravilock/goduit/internal/config/mongo"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	privateKeyFile, err := os.Open("./jwtRS256.key")
	if err != nil {
		log.Fatal("Failed to open private key file", err)
	}

	if err := encryptionkeys.LoadPrivateKey(privateKeyFile); err != nil {
		log.Fatal("Failed to load private key file content", err)
	}

	if err := privateKeyFile.Close(); err != nil {
		log.Fatal("Failed to close private key file", err)
	}

	publicKeyFile, err := os.Open("./jwtRS256.key.pub")
	if err != nil {
		log.Fatal("Failed to open public key file", err)
	}

	if err := encryptionkeys.LoadPublicKey(publicKeyFile); err != nil {
		log.Fatal("Failed to load public key file content", err)
	}

	if err := publicKeyFile.Close(); err != nil {
		log.Fatal("Failed to close publicKeyFile key file", err)
	}

	databaseURI := os.Getenv("DB_URL")
	if databaseURI == "" {
		log.Fatal("You must sey your 'DATABASE_URI' environmental variable.")
	}
	// Connect Mongo DB
	if err := mongo.ConnectDatabase(databaseURI); err != nil {
		log.Fatal("Error connecting to database", err)
	}
	defer mongo.DisconnectDatabase()

	// Echo instance
	e := echo.New()

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
	routers.ArticlesRouter(apiGroup)

	// Start server
	e.Logger.Fatal(e.Start(":6969"))
}
