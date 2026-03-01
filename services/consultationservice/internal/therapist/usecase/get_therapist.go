package usecase

import (
	"consultationservice/internal/therapist/domain"
	"consultationservice/internal/therapist/repository"
	"context"
)

type GetTherapistUseCase struct {
	repo repository.TherapistRepository
}

func NewGetTherapistUseCase(repo repository.TherapistRepository) *GetTherapistUseCase {
	return &GetTherapistUseCase{repo: repo}
}

func (uc *GetTherapistUseCase) Execute(ctx context.Context, profileID string) (*domain.Therapist, error) {
	return uc.repo.GetByProfileID(ctx, profileID)
}

func (uc *GetTherapistUseCase) ExecuteAll(ctx context.Context) ([]*domain.Therapist, error) {
	return uc.repo.GetAll(ctx)
}

func (uc *GetTherapistUseCase) Search(ctx context.Context, query string, limit, skip int64) ([]*domain.Therapist, error) {
	return uc.repo.Search(ctx, query, limit, skip)
}
