package rabbitmq

import (
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ravilock/goduit/internal/app"
)

type queueConsumer struct {
	handler      app.Handler
	logger       *slog.Logger
	connection   *amqp.Connection
	channel      *amqp.Channel
	queueName    string
	queue        amqp.Queue
	deliveryChan <-chan amqp.Delivery
	exitChan     chan struct{}
}

func NewQueueConsumer(handler app.Handler, connection *amqp.Connection, queueName string, logger *slog.Logger) (*queueConsumer, error) {
	logger = logger.With("emitter", "queue-consumer", "queue-name", queueName)
	queueConsumer := &queueConsumer{
		handler:    handler,
		logger:     logger,
		connection: connection,
		queueName:  queueName,
	}

	if err := queueConsumer.setupAmpq(); err != nil {
		return nil, err
	}

	return queueConsumer, nil
}

func (q *queueConsumer) setupAmpq() error {
	if err := q.setupAmpqChannel(); err != nil {
		return err
	}
	if err := q.setupAmpqQueue(); err != nil {
		return err
	}
	return q.setupAmqpConsumer()
}

func (q *queueConsumer) setupAmpqChannel() error {
	if q.channel != nil {
		return nil
	}

	channel, err := q.connection.Channel()
	if err != nil {
		return err
	}

	q.channel = channel
	return nil
}

func (q *queueConsumer) setupAmpqQueue() error {
	queue, err := q.channel.QueueDeclare(q.queueName, true, false, false, false, nil)
	if err != nil {
		q.channel = nil
		return nil
	}
	q.queue = queue
	return nil
}

func (q *queueConsumer) setupAmqpConsumer() error {
	deliveryChan, err := q.channel.Consume(q.queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}
	q.deliveryChan = deliveryChan
	return nil
}

func (q *queueConsumer) Consume() {
	for {
		q.logger.Debug("Waiting for new rabbitmq deliveries")
		select {
		case delivery := <-q.deliveryChan:
			q.logger.Info("Delivery received", "body", string(delivery.Body))
			message := NewAmpqMessage(&delivery)
			q.handler.Handle(message)
			q.logger.Info("Delivery sent")
		case <-q.exitChan:
			// TODO: close consumer and queue connection
			return
		}
	}
}

func (q *queueConsumer) Stop() {
	q.exitChan <- struct{}{}
}
