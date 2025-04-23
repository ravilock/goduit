package services

import (
	"context"
)

type userFollower interface {
	Follow(ctx context.Context, followed, following string) error
}

type FollowUserService struct {
	repository userFollower
}

func NewFollowUserService(repository userFollower) *FollowUserService {
	return &FollowUserService{
		repository: repository,
	}
}

// Follow establishes a follow relationship between two users.
//
// The followed parameter represents the ID of the user to be followed.
//
// The follower parameter represents the ID of the user that is following.
func (s *FollowUserService) Follow(ctx context.Context, followed, following string) error {
	return s.repository.Follow(ctx, followed, following)
}
