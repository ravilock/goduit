package models

type User struct {
	Username     *string `bson:"username"`
	Email        *string `bson:"email"`
	PasswordHash *string `bson:"passwordHash"`
	Bio          *string `bson:"bio,omitempty"`
	Image        *string `bson:"image,omitempty"`
}
