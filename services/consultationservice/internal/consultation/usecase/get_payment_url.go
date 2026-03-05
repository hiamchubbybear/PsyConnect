package usecase

import (
	"consultationservice/internal/consultation/repository"
	"context"
	"errors"
	"fmt"
)

type GetPaymentURLUseCase struct {
	sessionRepo repository.SessionRepository
}

func NewGetPaymentURLUseCase(repo repository.SessionRepository) *GetPaymentURLUseCase {
	return &GetPaymentURLUseCase{sessionRepo: repo}
}

func (uc *GetPaymentURLUseCase) Execute(ctx context.Context, sessionID string) (string, error) {
	session, err := uc.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return "", err
	}
	if session == nil {
		return "", errors.New("session not found")
	}

	
	return fmt.Sprintf("https://mock-payment-gateway.psyconnect.dev/pay?session_id=%s&amount=%.2f", sessionID, session.Price), nil
}
