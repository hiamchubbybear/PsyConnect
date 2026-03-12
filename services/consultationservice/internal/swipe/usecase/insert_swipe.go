package usecase

import (
	"consultationservice/internal/swipe/domain"
	"consultationservice/internal/swipe/repository"
	"context"
)

type InsertSwipeUseCase struct {
	swipeRepo repository.SwipeRepository
}

func NewInsertSwipeUseCase(swipeRepo repository.SwipeRepository) *InsertSwipeUseCase {
	return &InsertSwipeUseCase{
		swipeRepo: swipeRepo,
	}
}

func (uc *InsertSwipeUseCase) Execute(ctx context.Context, clientID, therapistID string, points float32, reasons []string, status string) error {
	swipe := domain.NewSwipe(clientID, therapistID, points, reasons)
	if status != "" {
		swipe.Status = status
	}
	return uc.swipeRepo.InsertSwipe(ctx, swipe)
}
