package services

import (
	"context"

	"github.com/ravilock/goduit/internal/app/repositories"
)

// IsFollowedBy determines wether or not a user follows another user. Returns bool
//
// The followed parameter represents the username of the user to be followed.
//
// The following parameter represents the username of the user that is following.
func IsFollowedBy(followed, following string, ctx context.Context) bool {
	if following == "" || followed == following {
		return false
	}
	_, err := repositories.IsFollowedBy(followed, following, ctx)
	return err == nil
}
