package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ravilock/goduit/api/validators"
	articleHandlers "github.com/ravilock/goduit/internal/articlePublisher/handlers"
	articleRepositories "github.com/ravilock/goduit/internal/articlePublisher/repositories"
	articleServices "github.com/ravilock/goduit/internal/articlePublisher/services"
	"github.com/ravilock/goduit/internal/config"
	"github.com/ravilock/goduit/internal/config/mongo"
	followerHandlers "github.com/ravilock/goduit/internal/followerCentral/handlers"
	followerRepositories "github.com/ravilock/goduit/internal/followerCentral/repositories"
	followerServices "github.com/ravilock/goduit/internal/followerCentral/services"
	"github.com/ravilock/goduit/internal/identity"
	profileHandlers "github.com/ravilock/goduit/internal/profileManager/handlers"
	profileRepositories "github.com/ravilock/goduit/internal/profileManager/repositories"
	profileServices "github.com/ravilock/goduit/internal/profileManager/services"
)

func main() {
	privateKeyFile, err := os.Open(os.Getenv("PRIVATE_KEY_LOCATION"))
	if err != nil {
		log.Fatal("Failed to open private key file", err)
	}

	if err := config.LoadPrivateKey(privateKeyFile); err != nil {
		log.Fatal("Failed to load private key file content", err)
	}

	if err := privateKeyFile.Close(); err != nil {
		log.Fatal("Failed to close private key file", err)
	}

	publicKeyFile, err := os.Open(os.Getenv("PUBLIC_KEY_LOCATION"))
	if err != nil {
		log.Fatal("Failed to open public key file", err)
	}

	if err := config.LoadPublicKey(publicKeyFile); err != nil {
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
	commentRepository := articleRepositories.NewCommentRepository(client)
	articlePublisherRepository := articleRepositories.NewArticleRepository(client)
	// services
	profileManager := profileServices.NewProfileManager(userRepository)
	followerCentral := followerServices.NewFollowerCentral(followerRepository)
	commentPublisher := articleServices.NewCommentPublisher(commentRepository)
	articlePublisher := articleServices.NewArticlePublisher(articlePublisherRepository)
	// handlers
	profileHandler := profileHandlers.NewProfileHandler(profileManager, followerCentral)
	followerHandler := followerHandlers.NewFollowerHandler(followerCentral, profileManager)
	articleHandler := articleHandlers.NewArticleHandler(articlePublisher, profileManager, followerCentral)
	commentHandler := articleHandlers.NewCommentHandler(commentPublisher, articlePublisher, profileManager, followerCentral)
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Start Validator
	if err := validators.InitValidator(); err != nil {
		log.Fatal("Could not init validator", err)
	}

	optionalAuthMiddleware := identity.CreateAuthMiddleware(false)
	requiredAuthMiddleware := identity.CreateAuthMiddleware(true)

	// Routes
	apiGroup := e.Group("/api")
	apiGroup.GET("/healthcheck", healthcheck)
	// User Routes
	usersGroup := apiGroup.Group("/users")
	usersGroup.POST("", profileHandler.Register)
	usersGroup.POST("/login", profileHandler.Login)
	userGroup := apiGroup.Group("/user")
	userGroup.GET("", profileHandler.GetOwnProfile, requiredAuthMiddleware)
	userGroup.PUT("", profileHandler.UpdateProfile, requiredAuthMiddleware)
	// Profile Routes
	profileGroup := apiGroup.Group("/profile")
	profileGroup.GET("/:username", profileHandler.GetProfile, optionalAuthMiddleware)
	profileGroup.POST("/:username/follow", followerHandler.Follow, requiredAuthMiddleware)
	profileGroup.DELETE("/:username/follow", followerHandler.Unfollow, requiredAuthMiddleware)
	// Article Routes
	articlesGroup := apiGroup.Group("/articles")
	articlesGroup.POST("", articleHandler.WriteArticle, requiredAuthMiddleware)
	articlesGroup.GET("", articleHandler.ListArticles, optionalAuthMiddleware)
	articleGroup := apiGroup.Group("/article")
	articleGroup.GET("/:slug", articleHandler.GetArticle, optionalAuthMiddleware)
	articleGroup.DELETE("/:slug", articleHandler.UnpublishArticle, requiredAuthMiddleware)
	articleGroup.PUT("/:slug", articleHandler.UpdateArticle, requiredAuthMiddleware)
	articleGroup.POST("/:slug/comments", commentHandler.WriteComment, requiredAuthMiddleware)
	articleGroup.GET("/:slug/comments", commentHandler.ListComments, optionalAuthMiddleware)
	articleGroup.DELETE("/:slug/comments/:id", commentHandler.DeleteComment, requiredAuthMiddleware)
	// Start server
	e.Logger.Fatal(e.Start(":6969"))
}

func healthcheck(c echo.Context) error {
	return c.String(http.StatusOK, fmt.Sprintln("OK"))
}
