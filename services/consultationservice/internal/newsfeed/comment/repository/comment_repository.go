package repository

import (
	"consultationservice/internal/newsfeed/comment/domain"
	"context"
)

type CommentRepository interface {
	
	CreateComment(ctx context.Context, comment *domain.Comment) error

	
	UpdateComment(ctx context.Context, id, content string) error

	
	DeleteComment(ctx context.Context, id string) error

	
	GetCommentByID(ctx context.Context, id string) (*domain.Comment, error)

	
	GetCommentsByPost(ctx context.Context, postID string, limit, skip int64) ([]domain.Comment, error)

	
	GetReplies(ctx context.Context, parentCommentID string) ([]domain.Comment, error)

	
	GetCommentsTree(ctx context.Context, postID string, maxDepth int, limit, skip int64) ([]domain.Comment, error)

	
	CountCommentsByPost(ctx context.Context, postID string) (int64, error)

	
	IncrementReplyCount(ctx context.Context, commentID string) error

	
	DecrementReplyCount(ctx context.Context, commentID string) error
}
