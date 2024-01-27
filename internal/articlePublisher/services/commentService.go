package services

import "github.com/ravilock/goduit/internal/articlePublisher/repositories"

type CommentPublisher struct {
	writeCommentService
}

func NewCommentPublisher(commentRepository *repositories.CommentRepository) *CommentPublisher {
	write := writeCommentService{commentRepository}
	return &CommentPublisher{write}
}
