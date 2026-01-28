package usecase

import (
	"consultationservice/internal/newsfeed/reaction/domain"
	"consultationservice/internal/newsfeed/reaction/repository"
	"context"
)

type GetReactionsUseCase struct {
	reactionRepo repository.ReactionRepository
}

func NewGetReactionsUseCase(reactionRepo repository.ReactionRepository) *GetReactionsUseCase {
	return &GetReactionsUseCase{
		reactionRepo: reactionRepo,
	}
}

func (uc *GetReactionsUseCase) Execute(ctx context.Context, postID string) ([]domain.Reaction, error) {
	return uc.reactionRepo.GetReactionsByPost(ctx, postID)
}
