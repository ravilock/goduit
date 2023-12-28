package repositories

import (
	"context"

	"github.com/ravilock/goduit/internal/app/models"
	db "github.com/ravilock/goduit/internal/config/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

func IsFollowedBy(from, to string, ctx context.Context) (*models.Follower, error) {
	var follower *models.Follower
	filter := bson.D{
		{Key: "from", Value: from},
		{Key: "to", Value: to},
	}
	collection := db.DatabaseClient.Database("conduit").Collection("followers")
	if err := collection.FindOne(ctx, filter).Decode(&follower); err != nil {
		return nil, err
	}
	return follower, nil
}
