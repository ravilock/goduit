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

type ArticleRepository struct {
	DBClient *mongo.Client
}

func NewArticleRepository(client *mongo.Client) *ArticleRepository {
	return &ArticleRepository{client}
}

func (r *ArticleRepository) WriteArticle(ctx context.Context, article *models.Article) error {
	now := time.Now().UTC().Truncate(time.Millisecond)
	article.CreatedAt = &now
	collection := r.DBClient.Database("conduit").Collection("articles")
	result, err := collection.InsertOne(ctx, article)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return app.ConflictError("articles")
		}
		return err
	}
	newId, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return errors.New("could not convert article ID")
	}
	article.ID = &newId
	return nil
}

func (r *ArticleRepository) ListArticles(ctx context.Context, author, tag string, limit, offset int64) ([]*models.Article, error) {
	filter := bson.D{}
	if author != "" {
		filter = append(filter, bson.E{Key: "author", Value: author})
	}
	if tag != "" {
		filter = append(filter, bson.E{
			Key: "tagList", Value: bson.D{{
				Key:   "$all",
				Value: []string{tag},
			}},
		})
	}
	opt := options.Find().SetLimit(limit).SetSkip(offset).SetSort(bson.D{{Key: "_id", Value: 1}})
	collection := r.DBClient.Database("conduit").Collection("articles")
	results := []*models.Article{}
	cursor, err := collection.Find(ctx, filter, opt)
	if err != nil {
		return results, err
	}
	if err = cursor.All(ctx, &results); err != nil {
		return results, err
	}
	return results, nil
}

func (r *ArticleRepository) GetArticleBySlug(ctx context.Context, slug string) (*models.Article, error) {
	var article *models.Article
	filter := bson.D{{
		Key:   "slug",
		Value: slug,
	}}
	collection := r.DBClient.Database("conduit").Collection("articles")
	if err := collection.FindOne(ctx, filter).Decode(&article); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, app.ArticleNotFoundError(slug, err)
		}
		return nil, err
	}
	return article, nil
}

func (r *ArticleRepository) DeleteArticle(ctx context.Context, slug string) error {
	filter := bson.D{{
		Key:   "slug",
		Value: slug,
	}}
	collection := r.DBClient.Database("conduit").Collection("articles")
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return app.ArticleNotFoundError(slug, nil)
	}
	return nil
}

func (r *ArticleRepository) UpdateArticle(ctx context.Context, slug string, article *models.Article) error {
	now := time.Now().UTC().Truncate(time.Millisecond)
	article.UpdatedAt = &now
	filter := bson.D{{Key: "slug", Value: slug}}
	update := bson.D{{Key: "$set", Value: article}}
	collection := r.DBClient.Database("conduit").Collection("articles")
	returnDocumentOption := options.After
	err := collection.FindOneAndUpdate(ctx, filter, update, &options.FindOneAndUpdateOptions{ReturnDocument: &returnDocumentOption}).Decode(article)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return app.ArticleNotFoundError(slug, err)
		}
		if mongo.IsDuplicateKeyError(err) {
			return app.ConflictError("articles")
		}
		return err
	}
	return nil
}
