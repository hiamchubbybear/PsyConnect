package repository

import (
	"consultationservice/internal/newsfeed/reaction/domain"
	"context"
)

type ReactionRepository interface {
	
	AddReaction(ctx context.Context, reaction *domain.Reaction) error

	
	RemoveReaction(ctx context.Context, postID, userID string) error

	
	GetUserReaction(ctx context.Context, postID, userID string) (*domain.Reaction, error)

	
	GetReactionsByPost(ctx context.Context, postID string) ([]domain.Reaction, error)

	
	CountReactionsByType(ctx context.Context, postID, reactionType string) (int64, error)
}
