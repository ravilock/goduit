package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Feed struct {
	ID         *primitive.ObjectID `bson:"_id,omitempty"`
	TargetUser *string             `bson:"userTarget"`
	Article    *string             `bson:"article"`
	CreatedAt  *time.Time          `bson:"createdAt"`
}
