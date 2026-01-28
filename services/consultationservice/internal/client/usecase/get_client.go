package usecase

import (
	"consultationservice/internal/client/domain"
	"consultationservice/internal/client/repository"
	"context"
)

type GetClientUseCase struct {
	repo repository.ClientRepository
}

func NewGetClientUseCase(repo repository.ClientRepository) *GetClientUseCase {
	return &GetClientUseCase{repo: repo}
}

func (uc *GetClientUseCase) Execute(ctx context.Context, profileID string) (*domain.Client, error) {
	return uc.repo.GetByProfileID(ctx, profileID)
}

func (uc *GetClientUseCase) ExecuteAll(ctx context.Context) ([]*domain.Client, error) {
	return uc.repo.GetAll(ctx)
}
