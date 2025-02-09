package api

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ravilock/goduit/api/validators"
	articleHandlers "github.com/ravilock/goduit/internal/articlePublisher/handlers"
	articlePublishers "github.com/ravilock/goduit/internal/articlePublisher/publishers"
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
	"github.com/ravilock/goduit/internal/rabbitmq"
	"github.com/spf13/viper"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

type Server interface {
	http.Handler
	Start()
}

type server struct {
	*echo.Echo
	db    *mongoDriver.Client
	queue *amqp.Connection
}

func (s *server) Start() {
	addr := fmt.Sprintf(":%d", viper.GetInt("port"))
	s.Logger.Fatal(s.Echo.Start(addr))
}

func NewServer() (Server, error) {
	serverLogger := log.NewLogger(map[string]string{"emitter": "Backstage-Groups-API"})
	databaseClient, err := mongo.ConnectDatabase(viper.GetString("db.url"))
	if err != nil {
		return nil, err
	}

	queueConnection, err := rabbitmq.ConnectQueue(viper.GetString("queue.url"))
	if err != nil {
		return nil, err
	}

	return createNewServer(databaseClient, queueConnection, serverLogger)
}

func createNewServer(databaseClient *mongoDriver.Client, queueConnection *amqp.Connection, _ *slog.Logger) (Server, error) {
	// TODO: Add logger to each controller
	// Echo instance
	e := echo.New()

	server := &server{
		Echo:  e,
		db:    databaseClient,
		queue: queueConnection,
	}
	// queue publishers
	articleQueuePublisher, err := articlePublishers.NewArticleQueuePublisher(queueConnection, viper.GetString("article.queue.name"))
	if err != nil {
		return nil, err
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
	articlePublisher := articleServices.NewArticlePublisher(articlePublisherRepository, articleQueuePublisher)
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
	articlesGroup.GET("/:slug", articleHandler.GetArticle, optionalAuthMiddleware)
	articlesGroup.DELETE("/:slug", articleHandler.UnpublishArticle, requiredAuthMiddleware)
	articlesGroup.PUT("/:slug", articleHandler.UpdateArticle, requiredAuthMiddleware)
	articlesGroup.POST("/:slug/comments", commentHandler.WriteComment, requiredAuthMiddleware)
	articlesGroup.GET("/:slug/comments", commentHandler.ListComments, optionalAuthMiddleware)
	articlesGroup.DELETE("/:slug/comments/:id", commentHandler.DeleteComment, requiredAuthMiddleware)
	return server, nil
}

func healthcheck(c echo.Context) error {
	return c.String(http.StatusOK, fmt.Sprintln("OK"))
}
