package usecase

import (
	"consultationservice/internal/newsfeed/reaction/domain"
	"consultationservice/internal/newsfeed/reaction/repository"
	"context"
)


type PostRepository interface {
	UpdateEngagementCount(ctx context.Context, id, field string, delta int) error
	GetPostByID(ctx context.Context, id string) (interface{}, error) 
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



func (uc *ToggleReactionUseCase) Execute(ctx context.Context, postID, userID, reactionType string) (*domain.Reaction, *domain.Reaction, error) {
	
	existing, err := uc.reactionRepo.GetUserReaction(ctx, postID, userID)
	if err != nil {
		return nil, nil, err
	}

	
	if existing != nil && existing.ReactionType == reactionType {
		err = uc.reactionRepo.RemoveReaction(ctx, postID, userID)
		if err != nil {
			return existing, nil, err
		}

		
		go uc.updateEngagementCount(postID, existing.ReactionType, -1)
		return existing, nil, nil
	}

	
	reaction := domain.NewReaction(postID, userID, reactionType)
	err = uc.reactionRepo.AddReaction(ctx, reaction)
	if err != nil {
		return existing, nil, err
	}

	
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
		
		uc.updateEngagementCount(postID, new.ReactionType, 1)
	} else if existing.ReactionType != new.ReactionType {
		
		uc.updateEngagementCount(postID, existing.ReactionType, -1)
		uc.updateEngagementCount(postID, new.ReactionType, 1)
	}
}
