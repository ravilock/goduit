package repositories

import (
	"context"

	"github.com/ravilock/goduit/internal/app/models"
	db "github.com/ravilock/goduit/internal/config/mongo"
)

func CreateArticle(article *models.Article, ctx context.Context) error {
	collection := db.DatabaseClient.Database("conduit").Collection("articles")
	if _, err := collection.InsertOne(ctx, article); err != nil {
		return err
	}
	return nil
}
