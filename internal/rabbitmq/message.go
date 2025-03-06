package rabbitmq

import (
	"github.com/ravilock/goduit/internal/app"

	amqp "github.com/rabbitmq/amqp091-go"
)

var _ app.Message = &ampqMessage{}

type ampqMessage struct {
	amqp.Delivery
}

func NewAmpqMessage(delivery amqp.Delivery) *ampqMessage {
	return &ampqMessage{
		Delivery: delivery,
	}
}

// Data implements app.Message.
func (a *ampqMessage) Data() []byte {
	return a.Delivery.Body
}

// Failure implements app.Message.
func (a *ampqMessage) Failure(error) error {
	return a.Reject(false)
}

// Success implements app.Message.
func (a *ampqMessage) Success() error {
	return a.Ack(false)
}
