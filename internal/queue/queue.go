package queue

import (
	"context"

	"github.com/ravilock/goduit/internal/app"
)

type Publisher interface {
	Publish(ctx context.Context, message []byte) error
	Close() error
}

type Consumer interface {
	Consume()
	Stop()
}

type Connection interface {
	NewPublisher(queueName string) (Publisher, error)
	NewConsumer(queueName string, handler app.Handler) (Consumer, error)
	Close() error
}
