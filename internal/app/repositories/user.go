package repositories

import (
	"context"

	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/internal/app/models"
	db "github.com/ravilock/goduit/internal/config/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterUser(user *models.User, ctx context.Context) error {
	collection := db.DatabaseClient.Database("conduit").Collection("users")
	if _, err := collection.InsertOne(ctx, user); err != nil {
		return err
	}
	return nil
}

func GetUserByEmail(email string, ctx context.Context) (*models.User, error) {
	var user *models.User
	filter := bson.D{{
		Key:   "email",
		Value: email,
	}}
	collection := db.DatabaseClient.Database("conduit").Collection("users")
	if err := collection.FindOne(ctx, filter).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, api.FailedLoginAttempt
		}
		return nil, err
	}
	return user, nil
}
