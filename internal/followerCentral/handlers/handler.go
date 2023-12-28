package handlers

import (
	"github.com/ravilock/goduit/internal/followerCentral/services"

	profileManager "github.com/ravilock/goduit/internal/profileManager/services"
)

type FollowerHandler struct {
	followUserHandler
	unfollowUserHandler
}

func NewFollowerHandler(central *services.FollowerCentral, manager *profileManager.ProfileManager) *FollowerHandler {
	follow := followUserHandler{central, manager}
	unfollow := unfollowUserHandler{central, manager}
	return &FollowerHandler{follow, unfollow}
}
