package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DatabaseClient *mongo.Client

func ConnectDatabase(databaseURI string) error {
	var err error
	DatabaseClient, err = mongo.Connect(context.Background(), options.Client().ApplyURI(databaseURI))
	if err != nil {
		return err
	}
	if err = testDatabase(); err != nil {
		return err
	}
	if err = ensureIndexes(); err != nil {
		return err
	}
	return nil
}

func DisconnectDatabase() {
	if err := DatabaseClient.Disconnect(context.Background()); err != nil {
		panic(err)
	}
}

func testDatabase() error {
	return DatabaseClient.Ping(context.Background(), nil)
}
