package repository

import (
	"consultationservice/internal/swipe/domain"
	"context"
)

type SwipeRepository interface {
	// InsertSwipe inserts a single swipe record
	InsertSwipe(ctx context.Context, swipe *domain.Swipe) error

	// InsertSwipes inserts multiple swipe records for a client
	InsertSwipes(ctx context.Context, clientID string, swipes []*domain.Swipe) error

	// SwipeAndMatch converts a swipe to a match
	// This deletes the swipe and creates a match record
	SwipeAndMatch(ctx context.Context, clientID, therapistID string, points float32, reasons []string) error

	// GetTopSwipes retrieves the top swipes for a client
	GetTopSwipes(ctx context.Context, clientID string, limit int) ([]*domain.Swipe, error)

	// DeleteSwipesByClient deletes all pending swipes for a client
	DeleteSwipesByClient(ctx context.Context, clientID string) error
}
