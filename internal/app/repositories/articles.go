package repositories

import (
	"context"

	"github.com/ravilock/goduit/internal/app"
	"github.com/ravilock/goduit/internal/app/models"
	db "github.com/ravilock/goduit/internal/config/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateArticle(article *models.Article, ctx context.Context) error {
	collection := db.DatabaseClient.Database("conduit").Collection("articles")
	if _, err := collection.InsertOne(ctx, article); err != nil {
		return err
	}
	return nil
}

func GetArticleBySlug(slug string, ctx context.Context) (*models.Article, error) {
	var article *models.Article
	filter := bson.D{{
		Key:   "slug",
		Value: slug,
	}}
	collection := db.DatabaseClient.Database("conduit").Collection("articles")
	if err := collection.FindOne(ctx, filter).Decode(&article); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, app.ArticleNotFoundError(slug, err)
		}
		return nil, err
	}
	return article, nil
}
