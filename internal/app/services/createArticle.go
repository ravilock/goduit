package services

import (
	"context"

	"github.com/ravilock/goduit/internal/app/models"
	"github.com/ravilock/goduit/internal/app/repositories"
)

func CreateArticle(model *models.Article, ctx context.Context) (*models.Article, error) {
	if err := repositories.CreateArticle(model, ctx); err != nil {
		return nil, err
	}

	return model, nil
}
