package services

import "github.com/ravilock/goduit/internal/articlePublisher/repositories"

type CommentPublisher struct {
	writeCommentService
	listCommentService
}

func NewCommentPublisher(commentRepository *repositories.CommentRepository) *CommentPublisher {
	write := writeCommentService{commentRepository}
	list := listCommentService{commentRepository}
	return &CommentPublisher{write, list}
}
