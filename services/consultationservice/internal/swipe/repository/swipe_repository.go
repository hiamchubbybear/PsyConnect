package repository

import (
	"consultationservice/internal/swipe/domain"
	"context"
)

type SwipeRepository interface {
	
	InsertSwipe(ctx context.Context, swipe *domain.Swipe) error

	
	InsertSwipes(ctx context.Context, clientID string, swipes []*domain.Swipe) error

	
	
	SwipeAndMatch(ctx context.Context, clientID, therapistID string, points float32, reasons []string) error

	
	GetTopSwipes(ctx context.Context, clientID string, limit int) ([]*domain.Swipe, error)

	
	DeleteSwipesByClient(ctx context.Context, clientID string) error
}
