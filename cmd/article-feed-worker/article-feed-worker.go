package main

import (
	"github.com/ravilock/goduit/internal/log"
	"github.com/ravilock/goduit/internal/mongo"

	articleRepositories "github.com/ravilock/goduit/internal/articlePublisher/repositories"
	articlefeed "github.com/ravilock/goduit/internal/articlePublisher/workers/article-feed"
	_ "github.com/ravilock/goduit/internal/config"
	followerRepositories "github.com/ravilock/goduit/internal/followerCentral/repositories"
	profileRepositories "github.com/ravilock/goduit/internal/profileManager/repositories"
	"github.com/ravilock/goduit/internal/rabbitmq"
	"github.com/spf13/viper"
)

func main() {
	logger := log.NewLogger(map[string]string{"emitter": "Goduit-Article-Feed-Worker"})
	forever := make(chan struct{})
	errChan := make(chan error)

	databaseClient, err := mongo.ConnectDatabase(viper.GetString("db.url"))
	if err != nil {
		panic(err)
	}

	userRepository := profileRepositories.NewUserRepository(databaseClient)
	followerRepository := followerRepositories.NewFollowerRepository(databaseClient)
	articlePublisherRepository := articleRepositories.NewArticleRepository(databaseClient)
	feedRepository := articleRepositories.NewFeedRepository(databaseClient)

	queueConnection, err := rabbitmq.ConnectQueue(viper.GetString("queue.url"))
	if err != nil {
		panic(err)
	}

	worker, err := articlefeed.NewArticleFeedWorker(queueConnection, articlePublisherRepository, userRepository, followerRepository, feedRepository, viper.GetString("article.queue.name"), errChan, logger)
	if err != nil {
		panic(err)
	}

	go worker.Consume()

	logger.Info(" [*] Waiting for messages. To exit press CTRL+C\n")
	select {
	case <-forever:
		return
	case err := <-errChan:
		logger.Error("Error during message consumption", "error", err.Error())
		return
	}
}
