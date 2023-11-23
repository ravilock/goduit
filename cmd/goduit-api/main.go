package main

import (
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ravilock/goduit/api/handlers"
	"github.com/ravilock/goduit/api/middlewares"
	"github.com/ravilock/goduit/api/validators"
	encryptionkeys "github.com/ravilock/goduit/internal/config/encryptionKeys"
	"github.com/ravilock/goduit/internal/config/mongo"
	followerHandlers "github.com/ravilock/goduit/internal/followerCentral/handlers"
	followerRepositories "github.com/ravilock/goduit/internal/followerCentral/repositories"
	followerServices "github.com/ravilock/goduit/internal/followerCentral/services"
	profileHandlers "github.com/ravilock/goduit/internal/profileManager/handlers"
	profileRepositories "github.com/ravilock/goduit/internal/profileManager/repositories"
	profileServices "github.com/ravilock/goduit/internal/profileManager/services"
)

func main() {
	privateKeyFile, err := os.Open(os.Getenv("PRIVATE_KEY_LOCATION"))
	if err != nil {
		log.Fatal("Failed to open private key file", err)
	}

	if err := encryptionkeys.LoadPrivateKey(privateKeyFile); err != nil {
		log.Fatal("Failed to load private key file content", err)
	}

	if err := privateKeyFile.Close(); err != nil {
		log.Fatal("Failed to close private key file", err)
	}

	publicKeyFile, err := os.Open(os.Getenv("PUBLIC_KEY_LOCATION"))
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
	client, err := mongo.ConnectDatabase(databaseURI)
	if err != nil {
		log.Fatal("Error connecting to database", err)
	}
	defer mongo.DisconnectDatabase(client)
	// repositories
	userRepository := profileRepositories.NewUserRepository(client)
	followerRepository := followerRepositories.NewFollowerRepository(client)
	// services
	profileManager := profileServices.NewProfileManager(userRepository)
	followerCentral := followerServices.NewFollowerCentral(followerRepository)
	// handlers
	profileHandler := profileHandlers.NewProfileHandler(profileManager, followerCentral)
	followerHandler := followerHandlers.NewFollowerHandler(followerCentral, profileManager)
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Start Validator
	if err := validators.InitValidator(); err != nil {
		log.Fatal("Could not init validator", err)
	}

	// Routes
	apiGroup := e.Group("/api")
	apiGroup.GET("/healthcheck", handlers.Healthcheck)
	// User Routes
	usersGroup := apiGroup.Group("/users")
	usersGroup.POST("", profileHandler.Register)
	usersGroup.POST("/login", profileHandler.Login)
	userGroup := apiGroup.Group("/user")
	userGroup.GET("", profileHandler.GetOwnProfile, middlewares.CreateAuthMiddleware(true))
	userGroup.PUT("", profileHandler.UpdateProfile, middlewares.CreateAuthMiddleware(true))
	// Profile Routes
	profileGroup := apiGroup.Group("/profile")
	profileGroup.GET("/:username", profileHandler.GetProfile, middlewares.CreateAuthMiddleware(false))
	profileGroup.POST("/:username/follow", followerHandler.Follow, middlewares.CreateAuthMiddleware(true))
	profileGroup.POST("/:username/unfollow", followerHandler.Unfollow, middlewares.CreateAuthMiddleware(true))
	// Article Routes
	articlesGroup := apiGroup.Group("/articles")
	articlesGroup.POST("", handlers.CreateArticle, middlewares.CreateAuthMiddleware(true))
	articleGroup := apiGroup.Group("/article")
	articleGroup.GET("/:slug", handlers.GetArticle, middlewares.CreateAuthMiddleware(false))
	// Start server
	e.Logger.Fatal(e.Start(":6969"))
}
