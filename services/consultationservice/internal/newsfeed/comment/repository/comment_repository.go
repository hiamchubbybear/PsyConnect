package repository

import (
	"consultationservice/internal/newsfeed/comment/domain"
	"context"
)

type CommentRepository interface {
	// CreateComment creates a new comment with depth/path calculation
	CreateComment(ctx context.Context, comment *domain.Comment) error

	// UpdateComment updates comment content
	UpdateComment(ctx context.Context, id, content string) error

	// DeleteComment soft-deletes a comment
	DeleteComment(ctx context.Context, id string) error

	// GetCommentByID retrieves a single comment
	GetCommentByID(ctx context.Context, id string) (*domain.Comment, error)

	// GetCommentsByPost gets top-level comments for a post
	GetCommentsByPost(ctx context.Context, postID string, limit, skip int64) ([]domain.Comment, error)

	// GetReplies gets direct replies to a comment
	GetReplies(ctx context.Context, parentCommentID string) ([]domain.Comment, error)

	// GetCommentsTree gets nested comment tree for a post
	GetCommentsTree(ctx context.Context, postID string, maxDepth int, limit, skip int64) ([]domain.Comment, error)

	// CountCommentsByPost counts all comments for a post
	CountCommentsByPost(ctx context.Context, postID string) (int64, error)

	// IncrementReplyCount increments reply count for a comment
	IncrementReplyCount(ctx context.Context, commentID string) error

	// DecrementReplyCount decrements reply count for a comment
	DecrementReplyCount(ctx context.Context, commentID string) error
}
