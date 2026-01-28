package usecase

import (
	"consultationservice/internal/consultation/repository"
	"context"
	"errors"
)

// DeleteSessionUseCase handles session deletion business logic
type DeleteSessionUseCase struct {
	sessionRepo repository.SessionRepository
}

// NewDeleteSessionUseCase creates a new DeleteSessionUseCase
func NewDeleteSessionUseCase(sessionRepo repository.SessionRepository) *DeleteSessionUseCase {
	return &DeleteSessionUseCase{
		sessionRepo: sessionRepo,
	}
}

// Execute deletes a session by ID
func (uc *DeleteSessionUseCase) Execute(ctx context.Context, sessionID string) error {
	if sessionID == "" {
		return errors.New("session ID is required")
	}

	// Check if session exists
	session, err := uc.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return err
	}

	// Check if session can be deleted (e.g., not already completed)
	if session.Status == "completed" {
		return errors.New("cannot delete completed session")
	}

	// Delete session
	return uc.sessionRepo.Delete(ctx, sessionID)
}
