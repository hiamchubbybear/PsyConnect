package usecase

import (
	"consultationservice/internal/matching/domain"
	"consultationservice/internal/matching/repository"
	"context"
)

type GetClientMatchesUseCase struct {
	matchRepo repository.MatchRepository
}

func NewGetClientMatchesUseCase(matchRepo repository.MatchRepository) *GetClientMatchesUseCase {
	return &GetClientMatchesUseCase{
		matchRepo: matchRepo,
	}
}

func (uc *GetClientMatchesUseCase) Execute(ctx context.Context, clientID string, page int64) ([]domain.Match, error) {
	return uc.matchRepo.GetMatchesByClientIDPaginated(ctx, clientID, page, 10)
}
