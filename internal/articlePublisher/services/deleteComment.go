package services

import (
	"context"
)

type commentDeleter interface {
	DeleteComment(ctx context.Context, ID string) error
}

type DeleteCommentService struct {
	repository commentDeleter
}

func NewDeleteCommentService(repository commentDeleter) *DeleteCommentService {
	return &DeleteCommentService{
		repository: repository,
	}
}

func (s *DeleteCommentService) DeleteComment(ctx context.Context, ID string) error {
	return s.repository.DeleteComment(ctx, ID)
}
