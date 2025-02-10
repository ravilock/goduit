package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Feed struct {
	UserID   *primitive.ObjectID `bson:"_id,omitempty"`
	Articles []FeedFragment      `bson:"articles,omitempty"`
}

type FeedFragment struct {
	ArticleID *string `bson:"articleID,omitempty"`
	Author    *string `bson:"author,omitempty"`
}
