package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ravilock/goduit/internal/app"
	"github.com/ravilock/goduit/internal/profileManager/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	DBClient *mongo.Client
}

func NewUserRepository(client *mongo.Client) *UserRepository {
	return &UserRepository{client}
}

func (r *UserRepository) RegisterUser(ctx context.Context, user *models.User) (*models.User, error) {
	now := time.Now().UTC().Truncate(time.Millisecond)
	user.CreatedAt = &now
	user.LastSession = &now
	collection := r.DBClient.Database("conduit").Collection("users")
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, app.ConflictError("users")
		}
		return nil, err
	}
	newId, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("Could not convert user ID")
	}
	user.ID = &newId
	return user, nil
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

func (r *UserRepository) GetUserByID(ctx context.Context, ID string) (*models.User, error) {
	var user *models.User
	userID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, fmt.Errorf("Could not parse ID: %s into ObjectID: %w", ID, err)
	}
	filter := bson.D{{
		Key:   "_id",
		Value: userID,
	}}
	collection := r.DBClient.Database("conduit").Collection("users")
	if err := collection.FindOne(ctx, filter).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, app.UserNotFoundError(ID, err)
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) UpdateProfile(ctx context.Context, subjectEmail, clientUsername string, user *models.User) error {
	now := time.Now().UTC().Truncate(time.Millisecond)
	user.UpdatedAt = &now
	filter := bson.D{
		{Key: "username", Value: clientUsername},
		{Key: "email", Value: subjectEmail},
	}
	update := bson.D{{Key: "$set", Value: user}}
	collection := r.DBClient.Database("conduit").Collection("users")
	returnDocumentOption := options.After
	err := collection.FindOneAndUpdate(ctx, filter, update, &options.FindOneAndUpdateOptions{ReturnDocument: &returnDocumentOption}).Decode(user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return app.ConflictError("users")
		}
		if err == mongo.ErrNoDocuments {
			return app.UserNotFoundError(fmt.Sprintf("%s+%s", subjectEmail, clientUsername), err)
		}
		return err
	}
	return nil
}
