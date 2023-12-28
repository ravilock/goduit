package services

import (
	"context"

	"github.com/ravilock/goduit/internal/app/repositories"
)

func IsFollowedBy(username, followerUsername string, ctx context.Context) bool {
	if followerUsername == "" {
		return false
	}
	_, err := repositories.IsFollowedBy(username, followerUsername, ctx)
	if err != nil {
		return false
	}
	return true
}
