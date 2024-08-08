package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DatabaseClient *mongo.Client

func ConnectDatabase(databaseURI string) (*mongo.Client, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(databaseURI))
	DatabaseClient = client
	if err != nil {
		return nil, err
	}
	if err = testDatabase(client); err != nil {
		return nil, err
	}
	if err = ensureIndexes(client); err != nil {
		return nil, err
	}
	return client, nil
}

func DisconnectDatabase(client *mongo.Client) {
	if err := client.Disconnect(context.Background()); err != nil {
		panic(err)
	}
}

func testDatabase(client *mongo.Client) error {
	return client.Ping(context.Background(), nil)
}
