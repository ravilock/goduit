package repositories

import (
	"context"

	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/internal/app"
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
			return nil, app.UserNotFoundError(email, err)
		}
		return nil, err
	}
	return user, nil
}

func GetUserByUsername(username string, ctx context.Context) (*models.User, error) {
	var user *models.User
	filter := bson.D{{
		Key:   "username",
		Value: username,
	}}
	collection := db.DatabaseClient.Database("conduit").Collection("users")
	if err := collection.FindOne(ctx, filter).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, app.UserNotFoundError(username, err)
		}
		return nil, err
	}
	return user, nil
}

func UpdateUser(user *models.User, ctx context.Context) (*models.User, error) {
	filter := bson.D{{Key: "email", Value: user.Email}}
	update := bson.D{{Key: "$set", Value: user}}
	collection := db.DatabaseClient.Database("conduit").Collection("users")
	updateResult, err := collection.UpdateOne(ctx, filter, update, nil)
	if err != nil {
		return nil, err
	}
	if updateResult.MatchedCount == 0 {
		return nil, api.UserNotFound(*user.Email)
	}
	return user, nil
}
