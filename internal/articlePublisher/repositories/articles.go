package repositories

import (
	"context"

	"github.com/ravilock/goduit/internal/app"
	"github.com/ravilock/goduit/internal/articlePublisher/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ArticleRepository struct {
	DBClient *mongo.Client
}

func NewArticleRepository(client *mongo.Client) *ArticleRepository {
	return &ArticleRepository{client}
}

func (r *ArticleRepository) WriteArticle(ctx context.Context, article *models.Article) error {
	collection := r.DBClient.Database("conduit").Collection("articles")
	if _, err := collection.InsertOne(ctx, article); err != nil {
		return err
	}
	return nil
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

func (r *ArticleRepository) UpdateArticle(ctx context.Context, slug string, article *models.Article) (*models.Article, error) {
	filter := bson.D{{Key: "slug", Value: slug}}
	update := bson.D{{Key: "$set", Value: article}}
	collection := r.DBClient.Database("conduit").Collection("articles")
	updateResult, err := collection.UpdateOne(ctx, filter, update, nil)
	if err != nil {
		return nil, err
	}
	if updateResult.MatchedCount == 0 {
		return nil, app.UserNotFoundError(slug, nil)
	}
	return article, nil
}
