package responses

import (
	profileManagerResponses "github.com/ravilock/goduit/internal/profileManager/responses"
	"time"
)

type ArticleResponse struct {
	Article Article `json:"article"`
}

type Article struct {
	Slug           string                          `json:"slug"`
	Title          string                          `json:"title"`
	Description    string                          `json:"description"`
	Body           string                          `json:"body"`
	TagList        []string                        `json:"tagList"`
	CreatedAt      *time.Time                      `json:"createdAt"`
	UpdatedAt      *time.Time                      `json:"updatedAt,omitempty"`
	Favorited      bool                            `json:"favorited"`
	FavoritesCount int64                           `json:"favoritesCount"`
	Author         profileManagerResponses.Profile `json:"author"`
}

type ArticlesResponse struct {
	Articles []Article `json:"articles"`
}
