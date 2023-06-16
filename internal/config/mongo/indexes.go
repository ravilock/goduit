package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ensureIndexes() {
	collection := DatabaseClient.Database("conduit").Collection("users")
	_, err := collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{"Username", 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		panic(err)
	}

	_, err = collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{"Email", 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		panic(err)
	}
}
