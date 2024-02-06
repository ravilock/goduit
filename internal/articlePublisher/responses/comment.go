package responses

import (
	profileManagerResponses "github.com/ravilock/goduit/internal/profileManager/responses"
	"time"
)

type CommentResponse struct {
	Comment Comment `json:"comment"`
}

type Comment struct {
	ID        string                          `json:"id"`
	CreatedAt time.Time                       `json:"createdAt"`
	UpdatedAt time.Time                       `json:"updatedAt,omitempty"`
	Body      string                          `json:"body"`
	Author    profileManagerResponses.Profile `json:"author"`
}

type CommentsResponse struct {
	Comment []Comment `json:"comments"`
}
