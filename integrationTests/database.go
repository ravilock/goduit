package integrationtests

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func ClearDatabase(client *mongo.Client) {
	conduitDb := client.Database("conduit")
	collections, err := conduitDb.ListCollectionNames(context.Background(), bson.D{})
	if err != nil {
		log.Fatal("Could not list collections", err)
	}
	for _, coll := range collections {
		_, err := conduitDb.Collection(coll).DeleteMany(context.Background(), bson.D{})
		if err != nil {
			log.Fatal("Could not clear database", err)
		}
	}
}
