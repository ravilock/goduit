package services

import (
	"context"

	"github.com/ravilock/goduit/internal/followerCentral/models"
)

type isFollowedChecker interface {
	IsFollowedBy(ctx context.Context, followed, following string) (*models.Follower, error)
}

type isFollowedByService struct {
	repository isFollowedChecker
}

// IsFollowedBy determines wether or not a user follows another user. Returns bool
//
// The followed parameter represents the ID of the user to be followed.
//
// The following parameter represents the ID of the user that is following.
func (s *isFollowedByService) IsFollowedBy(ctx context.Context, followed, following string) bool {
	if following == "" || followed == following {
		return false
	}
	_, err := s.repository.IsFollowedBy(ctx, followed, following)
	return err == nil
}
