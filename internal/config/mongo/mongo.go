package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DatabaseClient *mongo.Client

func ConnectDatabase(databaseURI string) {
	var err error
	DatabaseClient, err = mongo.Connect(context.Background(), options.Client().ApplyURI(databaseURI))
	if err != nil {
		panic(err)
	}
	testDatabase()
	ensureIndexes()
}

func DisconnectDatabase() {
	if err := DatabaseClient.Disconnect(context.Background()); err != nil {
		panic(err)
	}
}

func testDatabase() {
	if err := DatabaseClient.Ping(context.Background(), nil); err != nil {
		panic(err)
	}
}
