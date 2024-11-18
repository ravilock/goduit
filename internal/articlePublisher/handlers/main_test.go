package handlers

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ravilock/goduit/internal/articlePublisher/models"
	profileManagerModels "github.com/ravilock/goduit/internal/profileManager/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func assembleRandomUser() *profileManagerModels.User {
	ID := primitive.NewObjectID()
	username := uuid.NewString()
	email := username + "@test.test"
	image := "http://" + username + ".com"
	now := time.Now()
	return &profileManagerModels.User{
		ID:           &ID,
		Username:     &username,
		Email:        &email,
		PasswordHash: &username,
		Bio:          &username,
		Image:        &image,
		CreatedAt:    &now,
		UpdatedAt:    &now,
		LastSession:  &now,
	}
}

func assembleArticleModel(authorID primitive.ObjectID) *models.Article {
	articleID := primitive.NewObjectID()
	authorIDHex := authorID.Hex()
	articleTitle := "Article Title"
	articleSlug := "article-title"
	articleDescription := "Article Description"
	articleBody := "Article Body"
	articleTagList := []string{"test"}
	favoriteCounts := int64(0)
	now := time.Now()
	return &models.Article{
		ID:             &articleID,
		Author:         &authorIDHex,
		Slug:           &articleSlug,
		Title:          &articleTitle,
		Description:    &articleDescription,
		Body:           &articleBody,
		TagList:        articleTagList,
		CreatedAt:      &now,
		UpdatedAt:      &now,
		FavoritesCount: &favoriteCounts,
	}
}

func assembleCommentModel(authorID, articleID primitive.ObjectID, commentBody string) *models.Comment {
	commentID := primitive.NewObjectID()
	authorIDHex := authorID.Hex()
	articleIDHex := articleID.Hex()
	now := time.Now().UTC().Truncate(time.Millisecond)
	return &models.Comment{
		ID:        &commentID,
		Author:    &authorIDHex,
		Article:   &articleIDHex,
		Body:      &commentBody,
		CreatedAt: &now,
		UpdatedAt: &now,
	}
}

func makeSlug(title string) string {
	loweredTitle := strings.ToLower(title)
	return strings.ReplaceAll(loweredTitle, " ", "-")
}
