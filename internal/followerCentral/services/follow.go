package services

import "context"

type userFollower interface {
	Follow(ctx context.Context, followed, following string) error
}

type followUserService struct {
	repository userFollower
}

// Follow establishes a follow relationship between two users.
//
// The followed parameter represents the username of the user to be followed.
//
// The follower parameter represents the username of the user that is following.
func (s *followUserService) Follow(ctx context.Context, followed, following string) error {
	return s.repository.Follow(ctx, followed, following)
}
