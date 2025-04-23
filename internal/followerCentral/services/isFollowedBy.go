package services

import (
	"context"

	"github.com/ravilock/goduit/internal/followerCentral/models"
)

type isFollowedChecker interface {
	IsFollowedBy(ctx context.Context, followed, following string) (*models.Follower, error)
}

type IsFollowedByService struct {
	repository isFollowedChecker
}

func NewIsFollowedByService(repository isFollowedChecker) *IsFollowedByService {
	return &IsFollowedByService{
		repository: repository,
	}
}

// IsFollowedBy determines wether or not a user follows another user. Returns bool
//
// The followed parameter represents the ID of the user to be followed.
//
// The following parameter represents the ID of the user that is following.
func (s *IsFollowedByService) IsFollowedBy(ctx context.Context, followed, following string) bool {
	if following == "" || followed == following {
		return false
	}
	_, err := s.repository.IsFollowedBy(ctx, followed, following)
	return err == nil
}
