package app

type Message interface {
	Data() []byte
	Success() error
	Failure(error) error
}
