package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ravilock/goduit/internal/log"
	"github.com/ravilock/goduit/internal/mongo"

	articleRepositories "github.com/ravilock/goduit/internal/articlePublisher/repositories"
	articleFeedWorker "github.com/ravilock/goduit/internal/articlePublisher/workers/article-feed"
	_ "github.com/ravilock/goduit/internal/config"
	followerRepositories "github.com/ravilock/goduit/internal/followerCentral/repositories"
	profileRepositories "github.com/ravilock/goduit/internal/profileManager/repositories"
	"github.com/ravilock/goduit/internal/queue"
	"github.com/spf13/viper"
)

func main() {
	logger := log.NewLogger(map[string]string{"emitter": "Goduit-Article-Feed-Worker"})
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	databaseClient, err := mongo.ConnectDatabase(viper.GetString("db.url"))
	if err != nil {
		panic(err)
	}

	userRepository := profileRepositories.NewUserRepository(databaseClient)
	followerRepository := followerRepositories.NewFollowerRepository(databaseClient)
	articlePublisherRepository := articleRepositories.NewArticleRepository(databaseClient)
	feedRepository := articleRepositories.NewFeedRepository(databaseClient)

	queueConnection, err := queue.Connect(
		queue.QueueType(viper.GetString("queue.type")),
		viper.GetString("queue.url"),
	)
	if err != nil {
		logger.Error("Failed to connect to queue", "error", err)
		panic(err)
	}

	handler := articleFeedWorker.NewArticleFeedHandler(articlePublisherRepository, userRepository, followerRepository, feedRepository, logger)

	articleFeedQueueConsumer, err := queueConnection.NewConsumer(viper.GetString("article.queue.name"), handler)
	if err != nil {
		panic(err)
	}

	go articleFeedQueueConsumer.Consume()

	logger.Info(" [*] Waiting for messages. To exit press CTRL+C\n")
	<-sigChan
	articleFeedQueueConsumer.Stop()
}
