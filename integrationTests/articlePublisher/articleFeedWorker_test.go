package articlepublisher

import (
	"context"
	"log"
	"testing"
	"time"

	integrationtests "github.com/ravilock/goduit/integrationTests"
	"github.com/ravilock/goduit/internal/app"
	articleRepositories "github.com/ravilock/goduit/internal/articlePublisher/repositories"
	articleFeedWorker "github.com/ravilock/goduit/internal/articlePublisher/workers/article-feed"
	followerCentralRepositories "github.com/ravilock/goduit/internal/followerCentral/repositories"
	logger "github.com/ravilock/goduit/internal/log"
	"github.com/ravilock/goduit/internal/mongo"
	profileRepositories "github.com/ravilock/goduit/internal/profileManager/repositories"
	profileManagerRequests "github.com/ravilock/goduit/internal/profileManager/requests"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestArticleFeedHandler(t *testing.T) {
	// Connect Mongo DB
	client, err := mongo.ConnectDatabase(viper.GetString("db.url"))
	if err != nil {
		log.Fatalln("Error connecting to database", err)
	}
	followerCentralRepository := followerCentralRepositories.NewFollowerRepository(client)
	userRepository := profileRepositories.NewUserRepository(client)
	articlePublisherRepository := articleRepositories.NewArticleRepository(client)
	feedRepository := articleRepositories.NewFeedRepository(client)
	logger := logger.NewLogger(map[string]string{"emitter": "Goduit-Article-Feed-Worker"})
	messageMock := app.NewMockMessage(t)

	handler := articleFeedWorker.NewArticleFeedHandler(articlePublisherRepository, userRepository, followerCentralRepository, feedRepository, logger)

	t.Run("Should write followers feed when followed user posts new article", func(t *testing.T) {
		// Arrange
		authorIdentity, _ := integrationtests.MustRegisterUser(t, profileManagerRequests.RegisterPayload{})
		followerIdentity, followerToken := integrationtests.MustRegisterUser(t, profileManagerRequests.RegisterPayload{})
		integrationtests.MustFollowUser(t, authorIdentity.Username, followerToken)
		expectedArticle := integrationtests.GenerateArticleModel(authorIdentity.Subject)
		integrationtests.MustWriteArticleRegister(t, client, expectedArticle)
		messageData := []byte(expectedArticle.ID.Hex())
		messageMock.EXPECT().Data().Return(messageData)
		messageMock.EXPECT().Success().Return(nil)

		// Act
		handler.Handle(messageMock)

		// Assert
		feeds, err := feedRepository.PaginateFeed(context.Background(), followerIdentity.Subject, 1, 0)
		require.NoError(t, err)
		require.Len(t, feeds, 1)
		feed := feeds[0]
		require.Equal(t, expectedArticle.ID.Hex(), *feed.ArticleID)
		require.Equal(t, authorIdentity.Subject, *feed.Author)
	})
}

func TestArticleFeedWorkerWithQueue(t *testing.T) {
	// Connect Mongo DB
	client, err := mongo.ConnectDatabase(viper.GetString("db.url"))
	if err != nil {
		log.Fatalln("Error connecting to database", err)
	}
	followerCentralRepository := followerCentralRepositories.NewFollowerRepository(client)
	userRepository := profileRepositories.NewUserRepository(client)
	articlePublisherRepository := articleRepositories.NewArticleRepository(client)
	feedRepository := articleRepositories.NewFeedRepository(client)
	logger := logger.NewLogger(map[string]string{"emitter": "Goduit-Article-Feed-Worker-Integration-Test"})

	queueConnection := integrationtests.GetQueueConnection()

	t.Run("Should process article through queue and write to followers feed", func(t *testing.T) {
		// Arrange
		authorIdentity, _ := integrationtests.MustRegisterUser(t, profileManagerRequests.RegisterPayload{})
		followerIdentity, followerToken := integrationtests.MustRegisterUser(t, profileManagerRequests.RegisterPayload{})
		integrationtests.MustFollowUser(t, authorIdentity.Username, followerToken)
		expectedArticle := integrationtests.GenerateArticleModel(authorIdentity.Subject)
		integrationtests.MustWriteArticleRegister(t, client, expectedArticle)

		// Setup worker
		handler := articleFeedWorker.NewArticleFeedHandler(articlePublisherRepository, userRepository, followerCentralRepository, feedRepository, logger)
		consumer, err := queueConnection.NewConsumer(viper.GetString("article.queue.name")+"-integration-test", handler)
		require.NoError(t, err)

		// Start consumer
		go consumer.Consume()
		defer consumer.Stop()

		// Publish message
		publisher, err := queueConnection.NewPublisher(viper.GetString("article.queue.name") + "-integration-test")
		require.NoError(t, err)
		defer publisher.Close()

		err = publisher.Publish(context.Background(), []byte(expectedArticle.ID.Hex()))
		require.NoError(t, err)

		// Wait for message to be processed
		time.Sleep(2 * time.Second)

		// Assert
		feeds, err := feedRepository.PaginateFeed(context.Background(), followerIdentity.Subject, 1, 0)
		require.NoError(t, err)
		require.Len(t, feeds, 1)
		feed := feeds[0]
		require.Equal(t, expectedArticle.ID.Hex(), *feed.ArticleID)
		require.Equal(t, authorIdentity.Subject, *feed.Author)
	})
}
