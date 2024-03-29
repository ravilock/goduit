package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Follower represents a relationship between users, describing the followers that a particular user has.
//   - "Followed" represents the ID of the user to be followed
//   - "Follower" represents the ID of the user that is following
type Follower struct { // TODO: Think if this model should be named "Follow"
	ID       *primitive.ObjectID `bson:"_id,omitempty"`
	Followed *string             `bson:"followed"`
	Follower *string             `bson:"follower"`
}
