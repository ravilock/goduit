package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ensureIndexes(client *mongo.Client) error {
	usersCollection := client.Database("conduit").Collection("users")
	_, err := usersCollection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	_, err = usersCollection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	followersCollection := client.Database("conduit").Collection("followers")
	_, err = followersCollection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "followed", Value: 1},
			{Key: "follower", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	_, err = followersCollection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "follower", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	articlesCollection := client.Database("conduit").Collection("articles")
	_, err = articlesCollection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{{Key: "author", Value: 1}},
	})
	if err != nil {
		return err
	}

	_, err = articlesCollection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "slug", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	_, err = articlesCollection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{{Key: "tagList", Value: 1}},
	})
	if err != nil {
		return err
	}

	return nil
}
