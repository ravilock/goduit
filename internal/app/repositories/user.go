package repositories

import (
	"context"

	"github.com/ravilock/goduit/internal/app/models"
	db "github.com/ravilock/goduit/internal/config/mongo"
)

func RegisterUser(user *models.User, ctx context.Context) error {
	collection := db.DatabaseClient.Database("conduit").Collection("users")
	if _, err := collection.InsertOne(ctx, user); err != nil {
		return err
	}
	return nil
}
