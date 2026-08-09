package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Article struct {
	ID             *bson.ObjectID `bson:"_id,omitempty"`
	Author         *string        `bson:"author,omitempty"`
	Slug           *string        `bson:"slug,omitempty"`
	Title          *string        `bson:"title,omitempty"`
	Description    *string        `bson:"description,omitempty"`
	Body           *string        `bson:"body,omitempty"`
	TagList        []string       `bson:"tagList,omitempty"`
	CreatedAt      *time.Time     `bson:"createdAt,omitempty"`
	UpdatedAt      *time.Time     `bson:"updatedAt,omitempty"`
	FavoritesCount *int64         `bson:"favoritesCount,omitempty"`
}
