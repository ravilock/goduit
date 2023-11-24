package handlers

import (
	followerCentral "github.com/ravilock/goduit/internal/followerCentral/services"
	"github.com/ravilock/goduit/internal/profileManager/services"
)

type ProfileHandler struct {
	registerProfileHandler
	loginHandler
	updateProfileHandler
	getOwnProfileHandler
	getProfileHandler
}

func NewProfileHandler(manager *services.ProfileManager, central *followerCentral.FollowerCentral) *ProfileHandler {
	register := registerProfileHandler{manager}
	login := loginHandler{manager}
	updateProfile := updateProfileHandler{manager}
	getOwnProfile := getOwnProfileHandler{manager}
	getProfile := getProfileHandler{manager, central}
	return &ProfileHandler{register, login, updateProfile, getOwnProfile, getProfile}
}
