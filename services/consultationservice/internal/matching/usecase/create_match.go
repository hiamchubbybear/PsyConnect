package usecase

import (
	"context"
	"errors"

	clientRepo "consultationservice/internal/client/repository"
	"consultationservice/internal/matching/domain"
	"consultationservice/internal/matching/repository"
	therapistRepo "consultationservice/internal/therapist/repository"
)

type CreateMatchUseCase struct {
	matchRepo     repository.MatchRepository
	clientRepo    clientRepo.ClientRepository
	therapistRepo therapistRepo.TherapistRepository
}

func NewCreateMatchUseCase(
	matchRepo repository.MatchRepository,
	clientRepo clientRepo.ClientRepository,
	therapistRepo therapistRepo.TherapistRepository,
) *CreateMatchUseCase {
	return &CreateMatchUseCase{
		matchRepo:     matchRepo,
		clientRepo:    clientRepo,
		therapistRepo: therapistRepo,
	}
}

func (uc *CreateMatchUseCase) Execute(ctx context.Context, clientID, therapistID, source string, score float64, reasons []string) error {
	
	_, err := uc.clientRepo.GetByProfileID(ctx, clientID)
	if err != nil {
		return errors.New("client not found or invalid")
	}

	
	_, err = uc.therapistRepo.GetByProfileID(ctx, therapistID)
	
	if err != nil {
		return errors.New("therapist not found or invalid")
	}

	match := domain.NewMatch(clientID, therapistID, source, score, reasons)
	return uc.matchRepo.Create(ctx, match)
}
