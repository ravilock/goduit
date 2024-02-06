package services

import "github.com/ravilock/goduit/internal/articlePublisher/repositories"

type CommentPublisher struct {
	writeCommentService
	listCommentService
	getCommentService
	deleteCommentService
}

func NewCommentPublisher(commentRepository *repositories.CommentRepository) *CommentPublisher {
	write := writeCommentService{commentRepository}
	list := listCommentService{commentRepository}
	get := getCommentService{commentRepository}
	del := deleteCommentService{commentRepository}
	return &CommentPublisher{write, list, get, del}
}
