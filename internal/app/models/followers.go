package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Follower struct {
	ID   *primitive.ObjectID `bson:"_id"`
	From *string             `bson:"from"`
	To   *string             `bson:"to"`
}
