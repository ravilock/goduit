package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ravilock/goduit/internal/articlePublisher/models"
	"github.com/ravilock/goduit/internal/config/mongo"
	followerRepositories "github.com/ravilock/goduit/internal/followerCentral/repositories"
	followerServices "github.com/ravilock/goduit/internal/followerCentral/services"
)

func main() {
	databaseURI := os.Getenv("DB_URL")
	if databaseURI == "" {
		log.Fatalln("You must sey your 'DATABASE_URI' environmental variable.")
	}
	// Connect Mongo DB
	client, err := mongo.ConnectDatabase(databaseURI)
	if err != nil {
		log.Fatalln("Error connecting to database", err)
	}
	defer mongo.DisconnectDatabase(client)
	followerRepository := followerRepositories.NewFollowerRepository(client)
	followerCentral := followerServices.NewFollowerCentral(followerRepository)

	articleQueueName := os.Getenv("NEW_ARTICLE_QUEUE")
	if articleQueueName == "" {
		log.Fatalln("You must sey your 'NEW_ARTICLE_QUEUE' environmental variable.")
	}
	queueURI := os.Getenv("RABBIT_MQ_URL")
	if queueURI == "" {
		log.Fatalln("You must sey your 'RABBIT_MQ_URL' environmental variable.")
	}
	conn, err := amqp.Dial(queueURI)
	if err != nil {
		log.Fatalln("Failed to connect to RabbitMQ", err)
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalln("Failed to open RabbitMQ queue", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		articleQueueName, // name
		false,            // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Failed to declared queue %s", articleQueueName), err)
	}
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalln("Failed to register consumer", err)
	}
	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s\n", d.Body)
			article := new(models.Article)
			if err := json.Unmarshal(d.Body, article); err != nil {
				log.Printf("Could not unmarshal article: %s\n", err)
				continue
			}
			log.Printf("Received article data: %+v\n", article)
			followers, err := followerCentral.GetFollowers(context.Background(), *article.Author)
			if err != nil {
				log.Printf("Could not get article's author followers: %s\n", err)
				continue
			}
			log.Printf("Followers: %+v\n", followers)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
