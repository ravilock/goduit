package handlers

import (
	"github.com/ravilock/goduit/internal/cookie"
	followerCentral "github.com/ravilock/goduit/internal/followerCentral/services"
	"github.com/ravilock/goduit/internal/profileManager/services"
)

type ProfileHandler struct {
	registerProfileHandler
	loginHandler
	logoutHandler
	updateProfileHandler
	getOwnProfileHandler
	getProfileHandler
}

func NewProfileHandler(manager *services.ProfileManager, central *followerCentral.FollowerCentral, cookieManager *cookie.CookieManager) *ProfileHandler {
	register := registerProfileHandler{manager, cookieManager}
	login := loginHandler{manager, cookieManager}
	logout := logoutHandler{cookieManager}
	updateProfile := updateProfileHandler{manager, cookieManager}
	getOwnProfile := getOwnProfileHandler{manager}
	getProfile := getProfileHandler{manager, central}
	return &ProfileHandler{register, login, logout, updateProfile, getOwnProfile, getProfile}
}
