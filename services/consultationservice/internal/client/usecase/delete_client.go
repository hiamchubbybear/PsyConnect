package usecase

import (
	"consultationservice/internal/client/repository"
	"context"
)

type DeleteClientUseCase struct {
	repo repository.ClientRepository
}

func NewDeleteClientUseCase(repo repository.ClientRepository) *DeleteClientUseCase {
	return &DeleteClientUseCase{repo: repo}
}

func (uc *DeleteClientUseCase) Execute(ctx context.Context, profileID string) error {
	return uc.repo.Delete(ctx, profileID)
}
