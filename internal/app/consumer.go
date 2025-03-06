package app

type Consumer interface {
	Consume() <-chan Message
	StartConsumer()
	Stop()
}
