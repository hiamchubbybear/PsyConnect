package repository

import (
	"consultationservice/internal/newsfeed/post/domain"
	"context"
)

type PostRepository interface {
	CreatePost(ctx context.Context, post *domain.Post) error
	GetPostByID(ctx context.Context, id string) (*domain.Post, error)
	UpdatePost(ctx context.Context, id string, post *domain.Post) error
	DeletePost(ctx context.Context, id string) error
	GetPostsByUser(ctx context.Context, userID string, limit, skip int64) ([]domain.Post, error)
	GetFeed(ctx context.Context, userIDs []string, excludeIDs []string, limit, skip int64) ([]domain.Post, error)
	SearchPosts(ctx context.Context, query string, limit, skip int64) ([]domain.Post, error)
	GetTrendingPosts(ctx context.Context, limit int64) ([]domain.Post, error)
	GetPostsByTag(ctx context.Context, tag string, limit, skip int64) ([]domain.Post, error)
	GetPostsByCategory(ctx context.Context, category string, limit, skip int64) ([]domain.Post, error)
	IncrementViewCount(ctx context.Context, id string) error
	UpdateEngagementCount(ctx context.Context, id, field string, delta int) error
	GetPopularTags(ctx context.Context, limit int) ([]string, error)
}
