package queue

import (
	"context"
	"log/slog"
	"time"

	"github.com/ravilock/goduit/internal/app"
	"github.com/redis/go-redis/v9"
)

type redisConnection struct {
	client *redis.Client
}

func NewRedisConnection(url string) (Connection, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &redisConnection{client: client}, nil
}

func (r *redisConnection) NewPublisher(queueName string) (Publisher, error) {
	return &redisPublisher{
		client:    r.client,
		queueName: queueName,
	}, nil
}

func (r *redisConnection) NewConsumer(queueName string, handler app.Handler) (Consumer, error) {
	return newRedisConsumer(r.client, queueName, handler), nil
}

func (r *redisConnection) Close() error {
	if r.client != nil {
		return r.client.Close()
	}
	return nil
}

type redisPublisher struct {
	client    *redis.Client
	queueName string
}

func (p *redisPublisher) Publish(ctx context.Context, message []byte) error {
	return p.client.RPush(ctx, p.queueName, message).Err()
}

func (p *redisPublisher) Close() error {
	return nil
}

type redisConsumer struct {
	handler   app.Handler
	logger    *slog.Logger
	client    *redis.Client
	queueName string
	exitChan  chan struct{}
}

func newRedisConsumer(client *redis.Client, queueName string, handler app.Handler) *redisConsumer {
	logger := slog.Default().With("emitter", "redis-consumer", "queue-name", queueName)

	return &redisConsumer{
		handler:   handler,
		logger:    logger,
		client:    client,
		queueName: queueName,
		exitChan:  make(chan struct{}),
	}
}

func (c *redisConsumer) Consume() {
	ctx := context.Background()

	for {
		select {
		case <-c.exitChan:
			return
		default:
			c.logger.Debug("Waiting for new redis messages")

			result, err := c.client.BLPop(ctx, 1*time.Second, c.queueName).Result()
			if err != nil {
				if err == redis.Nil {
					continue
				}
				c.logger.Error("Error reading from redis queue", "error", err)
				continue
			}

			if len(result) < 2 {
				continue
			}

			messageData := result[1]
			c.logger.Info("Message received", "body", messageData)
			message := newRedisMessage([]byte(messageData))
			c.handler.Handle(message)
			c.logger.Info("Message processed")
		}
	}
}

func (c *redisConsumer) Stop() {
	close(c.exitChan)
}

type redisMessage struct {
	data []byte
}

func newRedisMessage(data []byte) app.Message {
	return &redisMessage{data: data}
}

func (m *redisMessage) Data() []byte {
	return m.data
}

func (m *redisMessage) Success() error {
	// Redis List operations (BLPOP) automatically remove messages from the queue.
	// Once consumed, the message is gone - there's no built-in acknowledgment.
	// This is a fire-and-forget pattern, unlike RabbitMQ's explicit ACK.
	return nil
}

func (m *redisMessage) Failure() error {
	// WARNING: Redis List operations don't support message requeuing.
	// Failed messages are lost. Consider using Redis Streams for at-least-once delivery.
	// For critical workloads, implement a dead-letter queue pattern manually.
	slog.Warn("Redis message failed but cannot be requeued - message will be lost",
		"data", string(m.data))
	return nil
}
