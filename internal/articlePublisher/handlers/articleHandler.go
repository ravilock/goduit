package handlers

import (
	"github.com/ravilock/goduit/internal/articlePublisher/producers"
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
}

func NewArticleHandler(
	publisher *services.ArticlePublisher,
	manager *profileManager.ProfileManager,
	central *followerCentral.FollowerCentral,
	producer *producers.ArticleProducer,
) *ArticleHandler {
	writeArticle := writeArticleHandler{publisher, manager, producer}
	getArticle := getArticleHandler{publisher, manager, central}
	unpublishArticle := unpublishArticleHandler{publisher}
	updateArticle := updateArticleHandler{publisher, manager}
	listArticles := listArticlesHandler{publisher, manager, central}
	return &ArticleHandler{writeArticle, getArticle, unpublishArticle, updateArticle, listArticles}
}
