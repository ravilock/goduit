package handlers

import (
	"github.com/ravilock/goduit/internal/articlePublisher/services"
	followerCentral "github.com/ravilock/goduit/internal/followerCentral/services"
	profileManager "github.com/ravilock/goduit/internal/profileManager/services"
)

type CommentHandler struct {
	writeCommentHandler
}

func NewCommentHandler(publisher *services.CommentPublisher, articlePublisher *services.ArticlePublisher, manager *profileManager.ProfileManager, central *followerCentral.FollowerCentral) *CommentHandler {
	writeComment := writeCommentHandler{publisher, articlePublisher, manager}
	return &CommentHandler{writeComment}
}