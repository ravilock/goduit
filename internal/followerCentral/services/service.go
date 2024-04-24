package services

import "github.com/ravilock/goduit/internal/followerCentral/repositories"

type FollowerCentral struct {
	followUserService
	unfollowUserService
	isFollowedByService
	getFollowersService
}

func NewFollowerCentral(followerRepository *repositories.FollowerRepository) *FollowerCentral {
	follow := followUserService{followerRepository}
	unfollow := unfollowUserService{followerRepository}
	isFollowedBy := isFollowedByService{followerRepository}
	getFollowers := getFollowersService{followerRepository}
	return &FollowerCentral{follow, unfollow, isFollowedBy, getFollowers}
}
