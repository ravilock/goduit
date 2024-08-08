package services

import (
	"context"
)

type commentDeleter interface {
	DeleteComment(ctx context.Context, ID string) error
}

type deleteCommentService struct {
	repository commentDeleter
}

func (s *deleteCommentService) DeleteComment(ctx context.Context, ID string) error {
	return s.repository.DeleteComment(ctx, ID)
}
