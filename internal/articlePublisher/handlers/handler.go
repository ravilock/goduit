package handlers

import (
	"github.com/ravilock/goduit/internal/articlePublisher/services"

	profileManager "github.com/ravilock/goduit/internal/profileManager/services"
)

type ArticleHandler struct {
	writeArticleHandler
}

func NewArticlehandler(publisher *services.ArticlePublisher, manager *profileManager.ProfileManager) *ArticleHandler {
	writeArticle := writeArticleHandler{publisher, manager}
	return &ArticleHandler{writeArticle}
}
