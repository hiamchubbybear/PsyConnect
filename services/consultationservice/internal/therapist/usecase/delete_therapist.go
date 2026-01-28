package usecase

import (
	"consultationservice/internal/therapist/repository"
	"context"
)

type DeleteTherapistUseCase struct {
	repo repository.TherapistRepository
}

func NewDeleteTherapistUseCase(repo repository.TherapistRepository) *DeleteTherapistUseCase {
	return &DeleteTherapistUseCase{repo: repo}
}

func (uc *DeleteTherapistUseCase) Execute(ctx context.Context, profileID string) error {
	return uc.repo.Delete(ctx, profileID)
}
