package services

import (
	"context"

	"github.com/ravilock/goduit/internal/app/models"
	"github.com/ravilock/goduit/internal/app/repositories"
)

func GetArticleBySlug(slug, username string, ctx context.Context) (*models.Article, error) {
	model, err := repositories.GetArticleBySlug(slug, ctx)
	if err != nil {
		return nil, err
	}
	return model, nil
}
