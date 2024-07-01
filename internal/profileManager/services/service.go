package services

import (
	"github.com/ravilock/goduit/internal/password"
	"github.com/ravilock/goduit/internal/profileManager/repositories"
)

type ProfileManager struct {
	registerProfileService
	logUserService
	updateProfileService
	getProfileService
}

func NewProfileManager(userRepository *repositories.UserRepository, hasherComparer password.HasherComparer) *ProfileManager {
	register := registerProfileService{userRepository, hasherComparer}
	login := logUserService{userRepository}
	updateProfile := updateProfileService{userRepository}
	getProfile := getProfileService{userRepository}
	return &ProfileManager{register, login, updateProfile, getProfile}
}
