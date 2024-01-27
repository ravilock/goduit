package services

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
)

type commentDeleter interface {
	DeleteComment(ctx context.Context, ID primitive.ObjectID) error
}

type deleteCommentService struct {
	repository commentDeleter
}

func (s *deleteCommentService) DeleteComment(ctx context.Context, ID primitive.ObjectID) error {
	return s.repository.DeleteComment(ctx, ID)
}
