package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Follower represents a relationship between users, describing the followers that a particular user has.
//   - "Followed" represents the username of the user to be followed
//   - "Follower" represents the username of the user that is following
type Follower struct { // TODO: Think if this model should be named "Follow"
	ID       *primitive.ObjectID `bson:"_id"`
	Followed *string             `bson:"followed"`
	Follower *string             `bson:"follower"`
}
