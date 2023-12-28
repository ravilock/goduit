package repositories

import (
	"context"

	"github.com/ravilock/goduit/internal/app/models"
	db "github.com/ravilock/goduit/internal/config/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

// Follow establishes a follow relationship between two users.
//
// The followed parameter represents the username of the user to be followed.
//
// The following parameter represents the username of the user that is following.
func Follow(followed, following string, ctx context.Context) error {
	follower := models.Follower{Following: &following, Followed: &followed}
	collection := db.DatabaseClient.Database("conduit").Collection("followers")
	if _, err := collection.InsertOne(ctx, follower); err != nil {
		return err
	}
	return nil
}

// IsFollowedBy queries for a follow relationship between two users. Returns *models.Follower.
//
// The followed parameter represents the username of the user to be followed.
//
// The following parameter represents the username of the user that is following.
func IsFollowedBy(followed, following string, ctx context.Context) (*models.Follower, error) {
	var follower *models.Follower
	filter := bson.D{
		{Key: "followed", Value: followed},
		{Key: "following", Value: following},
	}
	collection := db.DatabaseClient.Database("conduit").Collection("followers")
	if err := collection.FindOne(ctx, filter).Decode(&follower); err != nil {
		return nil, err
	}
	return follower, nil
}
