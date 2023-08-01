package services

import (
	"context"

	"github.com/ravilock/goduit/internal/app/dtos"
	"github.com/ravilock/goduit/internal/app/repositories"
	"github.com/ravilock/goduit/internal/app/transformers"
)

func GetArticleBySlug(slug, username string, ctx context.Context) (*dtos.Article, error) {
	model, err := repositories.GetArticleBySlug(slug, ctx)
	if err != nil {
		return nil, err
	}
	author, err := GetProfileByUsername(*model.Author, ctx)
	if err != nil {
		return nil, err
	}
	if username != "" {
		author.Following = IsFollowedBy(*author.Username, username, ctx)
	}
	return transformers.ArticleModelToDto(model, author), nil
}
