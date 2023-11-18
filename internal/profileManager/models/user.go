package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID           *primitive.ObjectID `bson:"_id,omitempty"`
	Username     *string             `bson:"username,omitempty"`
	Email        *string             `bson:"email,omitempty"`
	PasswordHash *string             `bson:"passwordHash,omitempty"`
	Bio          *string             `bson:"bio,omitempty"`
	Image        *string             `bson:"image,omitempty"`
}
