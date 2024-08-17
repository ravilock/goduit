package api

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ravilock/goduit/api/validators"
	articleHandlers "github.com/ravilock/goduit/internal/articlePublisher/handlers"
	articleRepositories "github.com/ravilock/goduit/internal/articlePublisher/repositories"
	articleServices "github.com/ravilock/goduit/internal/articlePublisher/services"
	followerHandlers "github.com/ravilock/goduit/internal/followerCentral/handlers"
	followerRepositories "github.com/ravilock/goduit/internal/followerCentral/repositories"
	followerServices "github.com/ravilock/goduit/internal/followerCentral/services"
	"github.com/ravilock/goduit/internal/identity"
	"github.com/ravilock/goduit/internal/log"
	"github.com/ravilock/goduit/internal/mongo"
	profileHandlers "github.com/ravilock/goduit/internal/profileManager/handlers"
	profileRepositories "github.com/ravilock/goduit/internal/profileManager/repositories"
	profileServices "github.com/ravilock/goduit/internal/profileManager/services"
	"github.com/spf13/viper"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

type Server interface {
	http.Handler
	Start()
}

type server struct {
	*echo.Echo
	db *mongoDriver.Client
}

func (s *server) Start() {
	s.Logger.Fatal(s.Echo.Start(viper.GetString(("server.address"))))
}

func NewServer() (Server, error) {
	serverLogger := log.NewLogger(map[string]string{"emitter": "Backstage-Groups-API"})
	databaseClient, err := mongo.ConnectDatabase(viper.GetString("db.url"))
	if err != nil {
		return nil, err
	}

	return createNewServer(databaseClient, serverLogger, false)
}

func createNewServer(databaseClient *mongoDriver.Client, logger *slog.Logger, testing bool) (Server, error) {
	// Echo instance
	e := echo.New()

	server := &server{
		Echo: e,
		db:   databaseClient,
		// measuresReporter: measuresReporter,
	}
	// repositories
	userRepository := profileRepositories.NewUserRepository(databaseClient)
	followerRepository := followerRepositories.NewFollowerRepository(databaseClient)
	commentRepository := articleRepositories.NewCommentRepository(databaseClient)
	articlePublisherRepository := articleRepositories.NewArticleRepository(databaseClient)
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

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Start Validator
	if err := validators.InitValidator(); err != nil {
		return nil, err
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
	profileGroup := apiGroup.Group("/profiles")
	profileGroup.GET("/:username", profileHandler.GetProfile, optionalAuthMiddleware)
	profileGroup.POST("/:username/followers", followerHandler.Follow, requiredAuthMiddleware)
	profileGroup.DELETE("/:username/followers", followerHandler.Unfollow, requiredAuthMiddleware)
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

	return server, nil
}

func healthcheck(c echo.Context) error {
	return c.String(http.StatusOK, fmt.Sprintln("OK"))
}
