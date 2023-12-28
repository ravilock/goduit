package services

import "context"

type userUnfollower interface {
	Unfollow(ctx context.Context, followed, following string) error
}

type unfollowUserService struct {
	repository userUnfollower
}

// Unfollow de-establishes a follow relationship between two users.
//
// The followed parameter represents the username of the user to be followed.
//
// The follower parameter represents the username of the user that is following.
func (s *unfollowUserService) Unfollow(ctx context.Context, followed, following string) error {
	return s.repository.Unfollow(ctx, followed, following)
}
