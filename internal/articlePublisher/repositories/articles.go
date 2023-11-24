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
