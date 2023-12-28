package services

import (
	"context"

	"github.com/ravilock/goduit/internal/app/repositories"
)

// Follow establishes a follow relationship between two users.
//
// The followed parameter represents the username of the user to be followed.
//
// The follower parameter represents the username of the user that is following.
func Follow(followed, following string, ctx context.Context) error {
	return repositories.Follow(followed, following, ctx)
}
