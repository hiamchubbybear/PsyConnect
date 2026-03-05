package usecase

import (
	"consultationservice/internal/consultation/domain"
	repo "consultationservice/internal/consultation/repository"
	"consultationservice/internal/kafka"
	"context"
	"errors"
	"fmt"
)

type StartCallUseCase struct {
	sessionRepo repo.SessionRepository
	producer    *kafka.Producer
}

func NewStartCallUseCase(sessionRepo repo.SessionRepository, producer *kafka.Producer) *StartCallUseCase {
	return &StartCallUseCase{
		sessionRepo: sessionRepo,
		producer:    producer,
	}
}

func (uc *StartCallUseCase) Execute(ctx context.Context, sessionID string, callerID string) error {
	
	session, err := uc.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	
	if !session.CanStartCall() {
		
		
	}

	
	
	var receiverID string
	if callerID == session.ClientID {
		receiverID = session.TherapistID
	} else if callerID == session.TherapistID {
		receiverID = session.ClientID
	} else {
		return errors.New("caller is not part of this session")
	}

	payload := domain.IncomingCallPayload{
		SessionID:  session.SessionID,
		CallerID:   callerID,
		CallerName: "Incoming Call", 
		
	}

	
	err = uc.producer.SendIncomingCallEvent(struct {
		ReceiverID string                     `json:"receiver_id"`
		Payload    domain.IncomingCallPayload `json:"payload"`
	}{
		ReceiverID: receiverID,
		Payload:    payload,
	})

	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}

	return nil
}
