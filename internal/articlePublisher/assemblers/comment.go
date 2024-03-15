package assemblers

import (
	"github.com/ravilock/goduit/internal/articlePublisher/models"
	"github.com/ravilock/goduit/internal/articlePublisher/responses"
	profileManagerResponses "github.com/ravilock/goduit/internal/profileManager/responses"
)

func CommentResponse(comment *models.Comment, author *profileManagerResponses.ProfileResponse) *responses.CommentResponse {
	response := new(responses.CommentResponse)
	response.Comment.ID = comment.ID.Hex()
	response.Comment.Body = *comment.Body
	response.Comment.CreatedAt = comment.CreatedAt
	response.Comment.UpdatedAt = comment.UpdatedAt
	response.Comment.Author = author.Profile
	return response
}
