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
	"github.com/ravilock/goduit/internal/rabbitmq"
	"github.com/spf13/viper"
)

func main() {
	logger := log.NewLogger(map[string]string{"emitter": "Goduit-Article-Feed-Worker"})
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

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

	articleFeedQueueConsumer, err := rabbitmq.NewQueueConsumer(queueConnection, viper.GetString("article.queue.name"))
	if err != nil {
		panic(err)
	}

	worker := articleFeedWorker.NewArticleFeedWorker(articleFeedQueueConsumer, articlePublisherRepository, userRepository, followerRepository, feedRepository, logger)

	go worker.Consume()

	logger.Info(" [*] Waiting for messages. To exit press CTRL+C\n")
	select {
	case <-sigChan:
		articleFeedQueueConsumer.Stop()
		return
	}
}
