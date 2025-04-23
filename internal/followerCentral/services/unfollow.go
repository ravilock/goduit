package services

import "context"

type userUnfollower interface {
	Unfollow(ctx context.Context, followed, following string) error
}

type UnfollowUserService struct {
	repository userUnfollower
}

func NewUnfollowUserService(repository userUnfollower) *UnfollowUserService {
	return &UnfollowUserService{
		repository: repository,
	}
}

// Unfollow de-establishes a follow relationship between two users.
//
// The followed parameter represents the ID of the user to be followed.
//
// The follower parameter represents the ID of the user that is following.
func (s *UnfollowUserService) Unfollow(ctx context.Context, followed, following string) error {
	return s.repository.Unfollow(ctx, followed, following)
}
