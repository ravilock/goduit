package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// FollowerModel represents a relationship between users
//   - "Followed" represents the username of the user to be followed
//   - "Follower" represents the username of the user that is following
type Follower struct {
	ID       *primitive.ObjectID `bson:"_id"`
	Followed *string             `bson:"followed"`
	Follower *string             `bson:"follower"`
}
