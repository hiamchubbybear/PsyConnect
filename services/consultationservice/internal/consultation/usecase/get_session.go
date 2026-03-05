package usecase

import (
	"consultationservice/internal/consultation/domain"
	"consultationservice/internal/consultation/repository"
	"context"
)


type GetSessionUseCase struct {
	sessionRepo repository.SessionRepository
}


func NewGetSessionUseCase(sessionRepo repository.SessionRepository) *GetSessionUseCase {
	return &GetSessionUseCase{
		sessionRepo: sessionRepo,
	}
}


func (uc *GetSessionUseCase) Execute(ctx context.Context, sessionID string) (*domain.Session, error) {
	return uc.sessionRepo.GetByID(ctx, sessionID)
}


func (uc *GetSessionUseCase) ExecuteByTherapist(ctx context.Context, therapistID string) ([]*domain.Session, error) {
	return uc.sessionRepo.GetByTherapist(ctx, therapistID)
}


func (uc *GetSessionUseCase) ExecuteByClient(ctx context.Context, clientID string) ([]*domain.Session, error) {
	return uc.sessionRepo.GetByClient(ctx, clientID)
}


func (uc *GetSessionUseCase) ExecuteAll(ctx context.Context) ([]*domain.Session, error) {
	return uc.sessionRepo.GetAll(ctx)
}
