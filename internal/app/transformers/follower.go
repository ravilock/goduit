package transformers

import (
	"github.com/ravilock/goduit/internal/app/dtos"
	"github.com/ravilock/goduit/internal/app/models"
)

func FollowerDtoToModel(follower *dtos.Follower) *models.Follower {
	return &models.Follower{
		From: follower.From,
		To:   follower.To,
	}
}
