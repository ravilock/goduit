package services

import "github.com/ravilock/goduit/internal/articlePublisher/repositories"

type CommentPublisher struct {
	writeCommentService
	listCommentService
	deleteCommentService
}

func NewCommentPublisher(commentRepository *repositories.CommentRepository) *CommentPublisher {
	write := writeCommentService{commentRepository}
	list := listCommentService{commentRepository}
	del := deleteCommentService{commentRepository}
	return &CommentPublisher{write, list, del}
}
