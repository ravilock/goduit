package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Article struct {
	ID             *primitive.ObjectID `bson:"_id,omitempty"`
	Author         *string             `bson:"author"`
	Slug           *string             `bson:"slug"`
	Title          *string             `bson:"title"`
	Description    *string             `bson:"description"`
	Body           *string             `bson:"body"`
	TagList        *[]string           `bson:"tagList"`
	CreatedAt      *time.Time          `bson:"createdAt"`
	UpdatedAt      *time.Time          `bson:"updatedAt"`
	FavoritesCount int64               `bson:"favoritesCount"`
}
