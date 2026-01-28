package repository

import (
	"consultationservice/internal/matching/domain"
	"context"
)

type MatchRepository interface {
	Create(ctx context.Context, match *domain.Match) error
	GetMatchesByClientID(ctx context.Context, clientID string) ([]domain.Match, error)
	GetMatchesByClientIDPaginated(ctx context.Context, clientID string, page int64, limit int64) ([]domain.Match, error)
}
