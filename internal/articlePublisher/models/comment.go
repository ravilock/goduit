package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID        *primitive.ObjectID `bson:"_id,omitempty"`
	Author    *string             `bson:"author,omitempty"`
	Article   *primitive.ObjectID `bson:"article,omitempty"`
	Body      *string             `bson:"body,omitempty"`
	CreatedAt *time.Time          `bson:"createdAt,omitempty"`
	UpdatedAt *time.Time          `bson:"updatedAt,omitempty"`
}
