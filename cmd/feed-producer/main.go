package main

import (
	"fmt"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	articleQueueName := os.Getenv("NEW_ARTICLE_QUEUE")
	if articleQueueName == "" {
		log.Fatal("You must sey your 'NEW_ARTICLE_QUEUE' environmental variable.")
	}
	queueURI := os.Getenv("RABBIT_MQ_URL")
	if queueURI == "" {
		log.Fatal("You must sey your 'RABBIT_MQ_URL' environmental variable.")
	}
	conn, err := amqp.Dial(queueURI)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ", err)
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open RabbitMQ queue", err)
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
		log.Fatal(fmt.Sprintf("Failed to declared queue %s", articleQueueName), err)
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
		log.Fatal("Failed to register consumer", err)
	}
	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
