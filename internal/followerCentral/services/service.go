package services

import "github.com/ravilock/goduit/internal/followerCentral/repositories"

type FollowerCentral struct {
	followUserService
	unfollowUserService
	isFollowedByService
}

func NewFollowerCentral(followerRepository *repositories.FollowerRepository) *FollowerCentral {
	follow := followUserService{followerRepository}
	unfollow := unfollowUserService{followerRepository}
	isFollowedBy := isFollowedByService{followerRepository}
	return &FollowerCentral{follow, unfollow, isFollowedBy}
}
