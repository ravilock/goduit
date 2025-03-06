package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ravilock/goduit/internal/app"
)

var _ app.Consumer = &queueConsumer{}

type queueConsumer struct {
	connection   *amqp.Connection
	channel      *amqp.Channel
	queueName    string
	queue        amqp.Queue
	deliveryChan <-chan amqp.Delivery
	messageChan  chan app.Message
	exitChan     chan struct{}
}

func NewQueueConsumer(connection *amqp.Connection, queueName string) (*queueConsumer, error) {
	queueConsumer := &queueConsumer{
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

func (q *queueConsumer) Consume() <-chan app.Message {
	return q.messageChan
}

func (q *queueConsumer) StartConsumer() {
	for {
		select {
		case message := <-q.deliveryChan:
			q.messageChan <- NewAmpqMessage(message)
		case <-q.exitChan:
			// TODO: close consumer and queue connection
			return
		}
	}
}

func (q *queueConsumer) Stop() {
	q.exitChan <- struct{}{}
}
