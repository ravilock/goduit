package handlers

import (
	"github.com/ravilock/goduit/internal/articlePublisher/services"

	followerCentral "github.com/ravilock/goduit/internal/followerCentral/services"
	profileManager "github.com/ravilock/goduit/internal/profileManager/services"
)

type ArticleHandler struct {
	writeArticleHandler
	getArticleHandler
	unpublishArticleHandler
	updateArticleHandler
	listArticlesHandler
	feedArticlesHandler
}

func NewArticleHandler(publisher *services.ArticlePublisher, manager *profileManager.ProfileManager, central *followerCentral.FollowerCentral) *ArticleHandler {
	writeArticle := writeArticleHandler{publisher, manager}
	getArticle := getArticleHandler{publisher, manager, central}
	unpublishArticle := unpublishArticleHandler{publisher}
	updateArticle := updateArticleHandler{publisher, manager}
	listArticles := listArticlesHandler{publisher, manager, central}
	feedArticles := feedArticlesHandler{publisher, manager}
	return &ArticleHandler{writeArticle, getArticle, unpublishArticle, updateArticle, listArticles, feedArticles}
}
