package queue

import (
	"fmt"
	"strings"
)

type QueueType string

const (
	RabbitMQ QueueType = "rabbitmq"
	Redis    QueueType = "redis"
)

func Connect(queueType QueueType, url string) (Connection, error) {
	switch strings.ToLower(string(queueType)) {
	case string(RabbitMQ):
		return NewRabbitMQConnection(url)
	case string(Redis):
		return NewRedisConnection(url)
	default:
		return nil, fmt.Errorf("unsupported queue type: %s", queueType)
	}
}
