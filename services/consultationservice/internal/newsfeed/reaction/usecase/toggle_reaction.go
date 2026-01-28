package usecase

import (
	"consultationservice/internal/newsfeed/reaction/domain"
	"consultationservice/internal/newsfeed/reaction/repository"
	"context"
)

// PostRepository interface for updating engagement counts
type PostRepository interface {
	UpdateEngagementCount(ctx context.Context, id, field string, delta int) error
	GetPostByID(ctx context.Context, id string) (interface{}, error) // Returns post for notification
}

type ToggleReactionUseCase struct {
	reactionRepo repository.ReactionRepository
	postRepo     PostRepository
}

func NewToggleReactionUseCase(
	reactionRepo repository.ReactionRepository,
	postRepo PostRepository,
) *ToggleReactionUseCase {
	return &ToggleReactionUseCase{
		reactionRepo: reactionRepo,
		postRepo:     postRepo,
	}
}

// Execute toggles a reaction (add, update, or remove)
// Returns: (existingReaction, newReaction, error)
func (uc *ToggleReactionUseCase) Execute(ctx context.Context, postID, userID, reactionType string) (*domain.Reaction, *domain.Reaction, error) {
	// Get existing reaction
	existing, err := uc.reactionRepo.GetUserReaction(ctx, postID, userID)
	if err != nil {
		return nil, nil, err
	}

	// If same reaction type, remove it (toggle off)
	if existing != nil && existing.ReactionType == reactionType {
		err = uc.reactionRepo.RemoveReaction(ctx, postID, userID)
		if err != nil {
			return existing, nil, err
		}

		// Update post engagement count
		go uc.updateEngagementCount(postID, existing.ReactionType, -1)
		return existing, nil, nil
	}

	// Add or update reaction
	reaction := domain.NewReaction(postID, userID, reactionType)
	err = uc.reactionRepo.AddReaction(ctx, reaction)
	if err != nil {
		return existing, nil, err
	}

	// Update post engagement counts
	go uc.updateEngagementCounts(postID, existing, reaction)

	return existing, reaction, nil
}

func (uc *ToggleReactionUseCase) updateEngagementCount(postID, reactionType string, delta int) {
	ctx := context.Background()
	field := "upvote_count"
	if reactionType == domain.VoteDown {
		field = "downvote_count"
	}
	uc.postRepo.UpdateEngagementCount(ctx, postID, field, delta)
}

func (uc *ToggleReactionUseCase) updateEngagementCounts(postID string, existing, new *domain.Reaction) {
	if existing == nil {
		// New reaction
		uc.updateEngagementCount(postID, new.ReactionType, 1)
	} else if existing.ReactionType != new.ReactionType {
		// Switched reaction
		uc.updateEngagementCount(postID, existing.ReactionType, -1)
		uc.updateEngagementCount(postID, new.ReactionType, 1)
	}
}
