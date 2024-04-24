package repositories

import (
	"context"

	"github.com/ravilock/goduit/internal/app"
	"github.com/ravilock/goduit/internal/followerCentral/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type FollowerRepository struct {
	DBClient *mongo.Client
}

func NewFollowerRepository(client *mongo.Client) *FollowerRepository {
	return &FollowerRepository{client}
}

// Follow establishes a follow relationship between two users.
//
// The followed parameter represents the ID of the user to be followed.
//
// The follower parameter represents the ID of the user that is following.
func (r *FollowerRepository) Follow(ctx context.Context, followed, follower string) error {
	followRelationship := models.Follower{Follower: &follower, Followed: &followed}
	collection := r.DBClient.Database("conduit").Collection("followers")
	if _, err := collection.InsertOne(ctx, followRelationship); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return app.ConflictError("followers")
		}
		return err
	}
	return nil
}

// Unfollow de-establishes a follow relationship between two users
//
// The followed parameter represents the ID of the user to be followed.
//
// The follower parameter represents the ID of the user that is following.
func (r *FollowerRepository) Unfollow(ctx context.Context, followed, follower string) error {
	filter := bson.D{
		{Key: "followed", Value: followed},
		{Key: "follower", Value: follower},
	}
	collection := r.DBClient.Database("conduit").Collection("followers")
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

// IsFollowedBy queries for a follow relationship between two users. Returns *models.Follower.
//
// The followed parameter represents the ID of the user to be followed.
//
// The follower parameter represents the ID of the user that is following.
func (r *FollowerRepository) IsFollowedBy(ctx context.Context, followed, follower string) (*models.Follower, error) {
	var followRelationship *models.Follower
	filter := bson.D{
		{Key: "followed", Value: followed},
		{Key: "follower", Value: follower},
	}
	collection := r.DBClient.Database("conduit").Collection("followers")
	if err := collection.FindOne(ctx, filter).Decode(&followRelationship); err != nil {
		return nil, err
	}
	return followRelationship, nil
}

// GetFollowers queries for all followers that a given user might have. Returns []*models.Follower.
//
// The followed parameter represents the ID of the user to be followed.
func (r *FollowerRepository) GetFollowers(ctx context.Context, followed string) ([]*models.Follower, error) { // Possibly needs pagination
	filter := bson.D{
		{Key: "followed", Value: followed},
	}
	collection := r.DBClient.Database("conduit").Collection("followers")
	results := []*models.Follower{}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}
