package usecase

import (
	"consultationservice/internal/consultation/repository"
	"context"
	"errors"
)


type DeleteSessionUseCase struct {
	sessionRepo repository.SessionRepository
}


func NewDeleteSessionUseCase(sessionRepo repository.SessionRepository) *DeleteSessionUseCase {
	return &DeleteSessionUseCase{
		sessionRepo: sessionRepo,
	}
}


func (uc *DeleteSessionUseCase) Execute(ctx context.Context, sessionID string) error {
	if sessionID == "" {
		return errors.New("session ID is required")
	}

	
	session, err := uc.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return err
	}

	
	if session.Status == "completed" {
		return errors.New("cannot delete completed session")
	}

	
	return uc.sessionRepo.Delete(ctx, sessionID)
}
