package usecase

import (
	"consultationservice/internal/consultation/repository"
	"consultationservice/internal/kafka"
	domain2 "consultationservice/internal/payment/domain"
	"context"
	"errors"

	"github.com/google/uuid"
)

type ProcessRefundUseCase struct {
	sessionRepo repository.SessionRepository
	producer    *kafka.Producer
}

func NewProcessRefundUseCase(repo repository.SessionRepository, producer *kafka.Producer) *ProcessRefundUseCase {
	return &ProcessRefundUseCase{
		sessionRepo: repo,
		producer:    producer,
	}
}

func (uc *ProcessRefundUseCase) Execute(ctx context.Context, sessionID string) (string, error) {
	session, err := uc.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return "", err
	}
	if session == nil {
		return "", errors.New("session not found")
	}

	if session.PaymentStatus != domain2.PaymentPaid {
		return "", errors.New("session is not paid, cannot refund")
	}

	
	refundTraceID := "REFUND-" + uuid.New().String()

	

	session.PaymentStatus = domain2.PaymentRefunded
	session.RefundTraceID = &refundTraceID
	session.Cancel("Refund processed")

	err = uc.sessionRepo.Update(ctx, session)
	if err != nil {
		return "", err
	}

	
	if uc.producer != nil {
		_ = uc.producer.SendSessionEvent("notification.push.consultation-updated", session)
	}

	return refundTraceID, nil
}
