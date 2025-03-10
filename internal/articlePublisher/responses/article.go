package responses

import (
	"time"

	profileManagerResponses "github.com/ravilock/goduit/internal/profileManager/responses"
)

type ArticleResponse struct {
	Article Article `json:"article"`
}

type Article struct {
	CreatedAt      *time.Time                      `json:"createdAt"`
	UpdatedAt      *time.Time                      `json:"updatedAt,omitempty"`
	Slug           string                          `json:"slug"`
	Title          string                          `json:"title"`
	Description    string                          `json:"description"`
	Body           string                          `json:"body"`
	Author         profileManagerResponses.Profile `json:"author"`
	TagList        []string                        `json:"tagList"`
	FavoritesCount int64                           `json:"favoritesCount"`
	Favorited      bool                            `json:"favorited"`
}

type ArticlesResponse struct {
	Articles []MultiArticle `json:"articles"`
}

type MultiArticle struct {
	CreatedAt      *time.Time                      `json:"createdAt"`
	UpdatedAt      *time.Time                      `json:"updatedAt,omitempty"`
	Slug           string                          `json:"slug"`
	Title          string                          `json:"title"`
	Description    string                          `json:"description"`
	Author         profileManagerResponses.Profile `json:"author"`
	TagList        []string                        `json:"tagList"`
	FavoritesCount int64                           `json:"favoritesCount"`
	Favorited      bool                            `json:"favorited"`
}
