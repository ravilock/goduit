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
	"github.com/ravilock/goduit/internal/cookie"
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
	serverLogger := log.NewLogger(map[string]string{"emitter": "Goduit-API"})
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
	feedRepository := articleRepositories.NewFeedRepository(databaseClient)
	// profile services
	registerProfileService := profileServices.NewRegisterProfileService(userRepository)
	logUserService := profileServices.NewLogUserService(userRepository)
	getProfileService := profileServices.NewGetProfileService(userRepository)
	updateUserService := profileServices.NewUpdateUserService(userRepository)
	// follower services
	followService := followerServices.NewFollowUserService(followerRepository)
	isFollowedByService := followerServices.NewIsFollowedByService(followerRepository)
	unfollowService := followerServices.NewUnfollowUserService(followerRepository)
	// comment services
	writeCommentService := articleServices.NewWriteCommentService(commentRepository)
	getCommentService := articleServices.NewGetCommentService(commentRepository)
	listCommentsService := articleServices.NewListCommentsService(commentRepository)
	deleteCommentService := articleServices.NewDeleteCommentService(commentRepository)
	// article services
	writeArticleService := articleServices.NewWriteArticleService(articlePublisherRepository, articleQueuePublisher)
	getArticleService := articleServices.NewGetArticleService(articlePublisherRepository)
	listArticlesService := articleServices.NewListArticlesService(articlePublisherRepository)
	feedArticlesService := articleServices.NewFeedArticlesService(articlePublisherRepository, feedRepository)
	updateArticleService := articleServices.NewUpdateArticleService(articlePublisherRepository)
	unpublishArticlesService := articleServices.NewUnpublishArticleService(articlePublisherRepository)
	// cookie manager
	cookieManager := cookie.NewCookieManager()
	// profile handlers
	registerProfileHandler := profileHandlers.NewRegisterProfileHandler(registerProfileService, cookieManager)
	getOwnProfileHandler := profileHandlers.NewGetOwnProfileHandler(getProfileService)
	getProfileHandler := profileHandlers.NewGetProfileHandler(getProfileService, isFollowedByService)
	loginHandler := profileHandlers.NewLoginHandler(logUserService, updateUserService, cookieManager)
	logoutHandler := profileHandlers.NewLogoutHandler(cookieManager)
	updateProfileHandler := profileHandlers.NewUpdateProfileHandler(updateUserService, cookieManager)
	// follower handlers
	followUserHandler := followerHandlers.NewFollowUserHandler(followService, getProfileService)
	unfollowUserHandler := followerHandlers.NewUnfollowUserHandler(unfollowService, getProfileService)
	// article handlers
	writeArticleHandler := articleHandlers.NewWriteArticleHandler(writeArticleService, getProfileService)
	getArticleHandler := articleHandlers.NewGetArticleHandler(getArticleService, getProfileService, isFollowedByService)
	listArticlesHandler := articleHandlers.NewListArticlesHandler(listArticlesService, getProfileService, isFollowedByService)
	feedArticlesHandler := articleHandlers.NewFeedArticlesHandler(feedArticlesService, getProfileService)
	updateArticleHandler := articleHandlers.NewUpdateArticleHandler(updateArticleService, getArticleService, getProfileService)
	unpublishArticlesHandler := articleHandlers.NewUnpublishArticleHandler(unpublishArticlesService, getArticleService)
	// comment handlers
	writeCommentHandler := articleHandlers.NewWriteCommentHandler(writeCommentService, getArticleService, getProfileService)
	listCommentsHandler := articleHandlers.NewListCommentsHandler(listCommentsService, getArticleService, getProfileService, isFollowedByService)
	deleteCommentHandler := articleHandlers.NewDeleteCommentHandler(deleteCommentService, getCommentService, getArticleService)

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// TODO: make origins be loaded as configurations
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))

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
	usersGroup.POST("", registerProfileHandler.Register)
	usersGroup.POST("/login", loginHandler.Login)
	usersGroup.POST("/logout", logoutHandler.Logout)
	userGroup := apiGroup.Group("/user")
	userGroup.GET("", getOwnProfileHandler.GetOwnProfile, requiredAuthMiddleware)
	userGroup.PUT("", updateProfileHandler.UpdateProfile, requiredAuthMiddleware)
	// Profile Routes
	profileGroup := apiGroup.Group("/profiles")
	profileGroup.GET("/:username", getProfileHandler.GetProfile, optionalAuthMiddleware)
	profileGroup.POST("/:username/followers", followUserHandler.Follow, requiredAuthMiddleware)
	profileGroup.DELETE("/:username/followers", unfollowUserHandler.Unfollow, requiredAuthMiddleware)
	// Article Routes
	articlesGroup := apiGroup.Group("/articles")
	articlesGroup.POST("", writeArticleHandler.WriteArticle, requiredAuthMiddleware)
	articlesGroup.GET("", listArticlesHandler.ListArticles, optionalAuthMiddleware)
	articlesGroup.GET("/feed", feedArticlesHandler.FeedArticles, requiredAuthMiddleware)
	articlesGroup.GET("/:slug", getArticleHandler.GetArticle, optionalAuthMiddleware)
	articlesGroup.DELETE("/:slug", unpublishArticlesHandler.UnpublishArticle, requiredAuthMiddleware)
	articlesGroup.PUT("/:slug", updateArticleHandler.UpdateArticle, requiredAuthMiddleware)
	articlesGroup.POST("/:slug/comments", writeCommentHandler.WriteComment, requiredAuthMiddleware)
	articlesGroup.GET("/:slug/comments", listCommentsHandler.ListComments, optionalAuthMiddleware)
	articlesGroup.DELETE("/:slug/comments/:id", deleteCommentHandler.DeleteComment, requiredAuthMiddleware)
	return server, nil
}

func healthcheck(c echo.Context) error {
	return c.String(http.StatusOK, fmt.Sprintln("OK"))
}
