package integrationtests

import (
	"log"

	"github.com/ravilock/goduit/internal/queue"
	"github.com/spf13/viper"
)

var queueConnection queue.Connection

func GetQueueConnection() queue.Connection {
	if queueConnection == nil {
		var err error
		queueConnection, err = queue.Connect(
			queue.QueueType(viper.GetString("queue.type")),
			viper.GetString("queue.url"),
		)
		if err != nil {
			log.Fatal("Failed to connect to queue:", err)
		}
	}
	return queueConnection
}

func CloseQueueConnection() {
	if queueConnection != nil {
		queueConnection.Close()
		queueConnection = nil
	}
}
