package repository

import (
	"consultationservice/internal/newsfeed/reaction/domain"
	"context"
)

type ReactionRepository interface {
	// AddReaction adds or updates a user's reaction to a post
	AddReaction(ctx context.Context, reaction *domain.Reaction) error

	// RemoveReaction removes a user's reaction from a post
	RemoveReaction(ctx context.Context, postID, userID string) error

	// GetUserReaction gets a specific user's reaction to a post
	GetUserReaction(ctx context.Context, postID, userID string) (*domain.Reaction, error)

	// GetReactionsByPost gets all reactions for a post
	GetReactionsByPost(ctx context.Context, postID string) ([]domain.Reaction, error)

	// CountReactionsByType counts reactions of a specific type for a post
	CountReactionsByType(ctx context.Context, postID, reactionType string) (int64, error)
}
