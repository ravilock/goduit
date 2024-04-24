package services

import (
	"context"

	"github.com/ravilock/goduit/internal/followerCentral/models"
)

type FollowersGetter interface {
	GetFollowers(ctx context.Context, followed string) ([]*models.Follower, error)
}

type getFollowersService struct {
	repository FollowersGetter
}

func (s *getFollowersService) GetFollowers(ctx context.Context, followed string) ([]string, error) {
	followRelationships, err := s.repository.GetFollowers(ctx, followed)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0, len(followRelationships))
	for _, followRelationship := range followRelationships {
		result = append(result, *followRelationship.Followed)
	}
	return result, nil
}
