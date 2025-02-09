package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

func ConnectQueue(url string) (*amqp.Connection, error) {
	return amqp.Dial(url)
}
