package repositories

import (
	"context"

	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/internal/app/models"
	db "github.com/ravilock/goduit/internal/config/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

// Follow establishes a follow relationship between two users.
//
// The followed parameter represents the username of the user to be followed.
//
// The follower parameter represents the username of the user that is following.
func Follow(followed, follower string, ctx context.Context) error {
	followRelationship := models.Follower{Follower: &follower, Followed: &followed}
	collection := db.DatabaseClient.Database("conduit").Collection("followers")
	if _, err := collection.InsertOne(ctx, followRelationship); err != nil {
		return err
	}
	return nil
}

// Unfollow de-establishes a follow relationship between two users
//
// The followed parameter represents the username of the user to be followed.
//
// The follower parameter represents the username of the user that is following.
func Unfollow(followed, follower string, ctx context.Context) error {
	filter := bson.D{
		{Key: "followed", Value: followed},
		{Key: "follower", Value: follower},
	}
	collection := db.DatabaseClient.Database("conduit").Collection("followers")
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return api.FollowerRelationshipNotFound(followed, follower)
	}
	return nil
}

// IsFollowedBy queries for a follow relationship between two users. Returns *models.Follower.
//
// The followed parameter represents the username of the user to be followed.
//
// The follower parameter represents the username of the user that is following.
func IsFollowedBy(followed, follower string, ctx context.Context) (*models.Follower, error) {
	var followRelationship *models.Follower
	filter := bson.D{
		{Key: "followed", Value: followed},
		{Key: "follower", Value: follower},
	}
	collection := db.DatabaseClient.Database("conduit").Collection("followers")
	if err := collection.FindOne(ctx, filter).Decode(&followRelationship); err != nil {
		return nil, err
	}
	return followRelationship, nil
}
