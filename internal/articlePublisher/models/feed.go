package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Feed struct {
	UserID   *bson.ObjectID `bson:"_id,omitempty"`
	Articles []FeedFragment `bson:"articles,omitempty"`
}

type FeedFragment struct {
	ArticleID *string `bson:"articleID,omitempty"`
	Author    *string `bson:"author,omitempty"`
}
