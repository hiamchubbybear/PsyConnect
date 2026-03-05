package usecase

import (
	"consultationservice/internal/consultation/repository"
	"consultationservice/internal/kafka"
	domain2 "consultationservice/internal/payment/domain"
	"context"
	"errors"
)

type ProcessPaymentWebhookUseCase struct {
	sessionRepo repository.SessionRepository
	producer    *kafka.Producer
}

func NewProcessPaymentWebhookUseCase(repo repository.SessionRepository, producer *kafka.Producer) *ProcessPaymentWebhookUseCase {
	return &ProcessPaymentWebhookUseCase{
		sessionRepo: repo,
		producer:    producer,
	}
}

type PaymentWebhookPayload struct {
	SessionID     string `json:"session_id"`
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"` 
}

func (uc *ProcessPaymentWebhookUseCase) Execute(ctx context.Context, payload PaymentWebhookPayload) error {
	session, err := uc.sessionRepo.GetByID(ctx, payload.SessionID)
	if err != nil {
		return err
	}
	if session == nil {
		return errors.New("session not found")
	}

	var newStatus domain2.PaymentStatus

	switch payload.Status {
	case "SUCCESS":
		newStatus = domain2.PaymentPaid
	case "FAILED":
		newStatus = domain2.PaymentFailed
	default:
		return errors.New("invalid payment status received in webhook")
	}

	
	err = uc.sessionRepo.UpdatePayment(ctx, payload.SessionID, newStatus, &payload.TransactionID)
	if err != nil {
		return err
	}

	
	if uc.producer != nil {
		_ = uc.producer.SendSessionEvent("notification.push.consultation-updated", session)
	}

	return nil
}
