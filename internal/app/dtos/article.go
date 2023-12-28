package dtos

import "time"

type Article struct {
	Slug           *string
	Title          *string
	Description    *string
	Body           *string
	TagList        *[]string
	CreatedAt      *time.Time
	UpdatedAt      *time.Time
	Favorited      bool
	FavoritesCount int64
	Author         *Profile
}
