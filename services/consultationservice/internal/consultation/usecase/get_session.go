package usecase

import (
	"consultationservice/internal/consultation/domain"
	"consultationservice/internal/consultation/repository"
	"context"
)

// GetSessionUseCase handles retrieving session information
type GetSessionUseCase struct {
	sessionRepo repository.SessionRepository
}

// NewGetSessionUseCase creates a new GetSessionUseCase
func NewGetSessionUseCase(sessionRepo repository.SessionRepository) *GetSessionUseCase {
	return &GetSessionUseCase{
		sessionRepo: sessionRepo,
	}
}

// Execute retrieves a session by ID
func (uc *GetSessionUseCase) Execute(ctx context.Context, sessionID string) (*domain.Session, error) {
	return uc.sessionRepo.GetByID(ctx, sessionID)
}

// ExecuteByTherapist retrieves all sessions for a therapist
func (uc *GetSessionUseCase) ExecuteByTherapist(ctx context.Context, therapistID string) ([]*domain.Session, error) {
	return uc.sessionRepo.GetByTherapist(ctx, therapistID)
}

// ExecuteByClient retrieves all sessions for a client
func (uc *GetSessionUseCase) ExecuteByClient(ctx context.Context, clientID string) ([]*domain.Session, error) {
	return uc.sessionRepo.GetByClient(ctx, clientID)
}

// ExecuteAll retrieves all sessions
func (uc *GetSessionUseCase) ExecuteAll(ctx context.Context) ([]*domain.Session, error) {
	return uc.sessionRepo.GetAll(ctx)
}
