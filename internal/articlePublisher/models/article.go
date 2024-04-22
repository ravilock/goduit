package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Article struct {
	ID             *primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Author         *string             `bson:"author,omitempty" json:"author,omitempty"`
	Slug           *string             `bson:"slug,omitempty" json:"slug,omitempty"`
	Title          *string             `bson:"title,omitempty" json:"title,omitempty"`
	Description    *string             `bson:"description,omitempty" json:"description,omitempty"`
	Body           *string             `bson:"body,omitempty" json:"body,omitempty"`
	TagList        []string            `bson:"tagList,omitempty" json:"tagList,omitempty"`
	CreatedAt      *time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt      *time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	FavoritesCount *int64              `bson:"favoritesCount,omitempty" json:"favoritesCount,omitempty"`
}
