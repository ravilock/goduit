package queue

import (
	"context"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ravilock/goduit/internal/app"
)

type rabbitmqConnection struct {
	conn *amqp.Connection
}

func NewRabbitMQConnection(url string) (Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	return &rabbitmqConnection{conn: conn}, nil
}

func (r *rabbitmqConnection) NewPublisher(queueName string) (Publisher, error) {
	return newRabbitMQPublisher(r.conn, queueName)
}

func (r *rabbitmqConnection) NewConsumer(queueName string, handler app.Handler) (Consumer, error) {
	return newRabbitMQConsumer(r.conn, queueName, handler)
}

func (r *rabbitmqConnection) Close() error {
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}

type rabbitmqPublisher struct {
	channel   *amqp.Channel
	queueName string
	queue     amqp.Queue
}

func newRabbitMQPublisher(conn *amqp.Connection, queueName string) (*rabbitmqPublisher, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	queue, err := channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		channel.Close()
		return nil, err
	}

	return &rabbitmqPublisher{
		channel:   channel,
		queueName: queueName,
		queue:     queue,
	}, nil
}

func (p *rabbitmqPublisher) Publish(ctx context.Context, message []byte) error {
	return p.channel.PublishWithContext(ctx, "", p.queue.Name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        message,
	})
}

func (p *rabbitmqPublisher) Close() error {
	if p.channel != nil {
		return p.channel.Close()
	}
	return nil
}

type rabbitmqConsumer struct {
	handler      app.Handler
	logger       *slog.Logger
	connection   *amqp.Connection
	channel      *amqp.Channel
	queueName    string
	queue        amqp.Queue
	deliveryChan <-chan amqp.Delivery
	exitChan     chan struct{}
}

func newRabbitMQConsumer(conn *amqp.Connection, queueName string, handler app.Handler) (*rabbitmqConsumer, error) {
	logger := slog.Default().With("emitter", "rabbitmq-consumer", "queue-name", queueName)
	
	consumer := &rabbitmqConsumer{
		handler:    handler,
		logger:     logger,
		connection: conn,
		queueName:  queueName,
		exitChan:   make(chan struct{}),
	}

	if err := consumer.setup(); err != nil {
		return nil, err
	}

	return consumer, nil
}

func (c *rabbitmqConsumer) setup() error {
	channel, err := c.connection.Channel()
	if err != nil {
		return err
	}
	c.channel = channel

	queue, err := c.channel.QueueDeclare(c.queueName, true, false, false, false, nil)
	if err != nil {
		c.channel.Close()
		return err
	}
	c.queue = queue

	deliveryChan, err := c.channel.Consume(c.queue.Name, "", false, false, false, false, nil)
	if err != nil {
		c.channel.Close()
		return err
	}
	c.deliveryChan = deliveryChan

	return nil
}

func (c *rabbitmqConsumer) Consume() {
	for {
		c.logger.Debug("Waiting for new rabbitmq deliveries")
		select {
		case delivery := <-c.deliveryChan:
			c.logger.Info("Delivery received", "body", string(delivery.Body))
			message := newRabbitMQMessage(&delivery)
			c.handler.Handle(message)
			c.logger.Info("Delivery sent")
		case <-c.exitChan:
			if c.channel != nil {
				c.channel.Close()
			}
			return
		}
	}
}

func (c *rabbitmqConsumer) Stop() {
	close(c.exitChan)
}

type rabbitmqMessage struct {
	*amqp.Delivery
}

func newRabbitMQMessage(delivery *amqp.Delivery) app.Message {
	return &rabbitmqMessage{Delivery: delivery}
}

func (m *rabbitmqMessage) Data() []byte {
	return m.Body
}

func (m *rabbitmqMessage) Success() error {
	return m.Ack(false)
}

func (m *rabbitmqMessage) Failure() error {
	return m.Reject(false)
}
