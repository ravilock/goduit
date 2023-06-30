package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	Username     *string            `bson:"username"`
	Email        *string            `bson:"email"`
	PasswordHash *string            `bson:"passwordHash"`
	Bio          *string            `bson:"bio,omitempty"`
	Image        *string            `bson:"image,omitempty"`
}
