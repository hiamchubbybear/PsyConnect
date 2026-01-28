package usecase

import (
	"consultationservice/internal/consultation/domain"
	repo "consultationservice/internal/consultation/repository"
	payment "consultationservice/internal/payment/domain"
	"context"
	"errors"
)

type ActivateSessionRequest struct {
	SessionID string
	ActorID   string
}

type ActivateSessionUseCase struct {
	repo repo.SessionRepository
}

func NewActivateSessionUseCase(sessionRepo repo.SessionRepository) *ActivateSessionUseCase {
	return &ActivateSessionUseCase{repo: sessionRepo}
}

func (uc *ActivateSessionUseCase) Execute(
	ctx context.Context,
	req ActivateSessionRequest,
) error {

	session, err := uc.repo.GetByID(ctx, req.SessionID)
	if err != nil {
		return err
	}

	if session.Status != domain.SessionStatusPending {
		return errors.New("session is not in pending state")
	}

	if session.PaymentStatus != payment.PaymentPaid {
		return errors.New("payment not completed")
	}

	session.Activate()

	return uc.repo.UpdateStatus(
		ctx,
		session.SessionID,
		session.Status,
	)
}
