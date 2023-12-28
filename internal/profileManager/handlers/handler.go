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
}

func NewProfileHandler(manager *services.ProfileManager, central *followerCentral.FollowerCentral) *ProfileHandler {
	register := registerProfileHandler{manager}
	login := loginHandler{manager}
	updateProfile := updateProfileHandler{manager}
	getOwnProfile := getOwnProfileHandler{manager}
	return &ProfileHandler{register, login, updateProfile, getOwnProfile}
}
