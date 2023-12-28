package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ensureIndexes() error {
	usersCollection := DatabaseClient.Database("conduit").Collection("users")
	_, err := usersCollection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{"username", 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	_, err = usersCollection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{"email", 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	followersCollection := DatabaseClient.Database("conduit").Collection("followers")
	_, err = followersCollection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{"followed", 1},
			{"follower", 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	_, err = followersCollection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{"follower", 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	articlesCollection := DatabaseClient.Database("conduit").Collection("articles")
	_, err = articlesCollection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{{"author", 1}},
	})
	if err != nil {
		return err
	}

	_, err = articlesCollection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{"slug", 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	_, err = articlesCollection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{{"tagList", 1}},
	})
	if err != nil {
		return err
	}

	return nil
}
