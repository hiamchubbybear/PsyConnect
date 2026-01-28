package repository

import (
	"consultationservice/internal/newsfeed/social/domain"
	"context"
)

type FollowRepository interface {
	FollowUser(ctx context.Context, followerID, followingID string) error
	UnfollowUser(ctx context.Context, followerID, followingID string) error
	GetFollowers(ctx context.Context, userID string) ([]domain.Follow, error)
	GetFollowing(ctx context.Context, userID string) ([]domain.Follow, error)
	IsFollowing(ctx context.Context, followerID, followingID string) (bool, error)
	CountFollowers(ctx context.Context, userID string) (int64, error)
	CountFollowing(ctx context.Context, userID string) (int64, error)
}

type BookmarkRepository interface {
	CreateBookmark(ctx context.Context, userID, postID string) error
	DeleteBookmark(ctx context.Context, userID, postID string) error
	GetBookmarks(ctx context.Context, userID string, limit, skip int64) ([]domain.Bookmark, error)
	IsBookmarked(ctx context.Context, userID, postID string) (bool, error)
}
