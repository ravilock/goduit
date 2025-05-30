package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/ravilock/goduit/internal/app"
	"github.com/ravilock/goduit/internal/articlePublisher/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CommentRepository struct {
	DBClient *mongo.Client
}

func NewCommentRepository(client *mongo.Client) *CommentRepository {
	return &CommentRepository{client}
}

func (r *CommentRepository) WriteComment(ctx context.Context, comment *models.Comment) error {
	now := time.Now().UTC().Truncate(time.Millisecond)
	comment.CreatedAt = &now
	collection := r.DBClient.Database("conduit").Collection("comments")
	result, err := collection.InsertOne(ctx, comment)
	if err != nil {
		return err
	}
	newId, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return errors.New("could not convert comment ID")
	}
	comment.ID = &newId
	return nil
}

func (r *CommentRepository) ListComments(ctx context.Context, article string) ([]*models.Comment, error) {
	filter := bson.D{{
		Key:   "article",
		Value: article,
	}}
	opt := options.Find().SetSort(bson.D{{Key: "_id", Value: -1}})
	collection := r.DBClient.Database("conduit").Collection("comments")
	results := []*models.Comment{}
	cursor, err := collection.Find(ctx, filter, opt)
	if err != nil {
		return results, err
	}
	if err = cursor.All(ctx, &results); err != nil {
		return results, err
	}
	return results, nil
}

func (r *CommentRepository) GetCommentByID(ctx context.Context, ID string) (*models.Comment, error) {
	var comment *models.Comment
	commentID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{
		Key:   "_id",
		Value: commentID,
	}}
	collection := r.DBClient.Database("conduit").Collection("comments")
	if err := collection.FindOne(ctx, filter).Decode(&comment); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, app.CommentNotFoundError(ID, err)
		}
		return nil, err
	}
	return comment, nil
}

func (r *CommentRepository) DeleteComment(ctx context.Context, ID string) error {
	commentID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return err
	}
	filter := bson.D{{
		Key:   "_id",
		Value: commentID,
	}}
	collection := r.DBClient.Database("conduit").Collection("comments")
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return app.CommentNotFoundError(ID, nil)
	}
	return nil
}
