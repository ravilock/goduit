package repositories

import (
	"context"

	"github.com/ravilock/goduit/internal/app"
	"github.com/ravilock/goduit/internal/profileManager/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	DBClient *mongo.Client
}

func NewUserRepository(client *mongo.Client) *UserRepository {
	return &UserRepository{client}
}

func (r *UserRepository) RegisterUser(ctx context.Context, user *models.User) error {
	collection := r.DBClient.Database("conduit").Collection("users")
	if _, err := collection.InsertOne(ctx, user); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return app.ConflictError("users")
		}
		return err
	}
	return nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user *models.User
	filter := bson.D{{
		Key:   "email",
		Value: email,
	}}
	collection := r.DBClient.Database("conduit").Collection("users")
	if err := collection.FindOne(ctx, filter).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, app.UserNotFoundError(email, err)
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user *models.User
	filter := bson.D{{
		Key:   "username",
		Value: username,
	}}
	collection := r.DBClient.Database("conduit").Collection("users")
	if err := collection.FindOne(ctx, filter).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, app.UserNotFoundError(username, err)
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) UpdateProfile(ctx context.Context, subjectEmail, clientUsername string, user *models.User) (*models.User, error) {
	filter := bson.D{
		{Key: "username", Value: clientUsername},
		{Key: "email", Value: subjectEmail},
	}
	update := bson.D{{Key: "$set", Value: user}}
	collection := r.DBClient.Database("conduit").Collection("users")
	updateResult, err := collection.UpdateOne(ctx, filter, update, nil)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, app.ConflictError("users")
		}
		return nil, err
	}
	if updateResult.MatchedCount == 0 {
		return nil, app.UserNotFoundError(*user.Email, nil)
	}
	return user, nil
}
