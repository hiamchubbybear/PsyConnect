package usecase

import (
	"consultationservice/internal/swipe/repository"
	"context"
)

type SwipeAndMatchUseCase struct {
	swipeRepo repository.SwipeRepository
}

func NewSwipeAndMatchUseCase(swipeRepo repository.SwipeRepository) *SwipeAndMatchUseCase {
	return &SwipeAndMatchUseCase{
		swipeRepo: swipeRepo,
	}
}

func (uc *SwipeAndMatchUseCase) Execute(ctx context.Context, clientID, therapistID string, points float32, reasons []string) error {
	return uc.swipeRepo.SwipeAndMatch(ctx, clientID, therapistID, points, reasons)
}
