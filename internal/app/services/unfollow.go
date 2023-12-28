package services

import (
	"context"

	"github.com/ravilock/goduit/internal/app/repositories"
)

// Unfollow de-establishes a follow relationship between two users.
//
// The followed parameter represents the username of the user to be followed.
//
// The follower parameter represents the username of the user that is following.
func Unfollow(followed, following string, ctx context.Context) error {
	return repositories.Unfollow(followed, following, ctx)
}
