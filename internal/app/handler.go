package app

type Handler interface {
	Handle(message Message)
}
