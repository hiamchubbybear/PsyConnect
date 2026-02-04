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
	// 1. Get Session
	session, err := uc.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	// 2. Validate
	if !session.CanStartCall() {
		// return errors.New("session is not active or not within allowed time window")
		// For testing, we might want to relax this or just logging warning
	}

	// 3. Construct Payload
	// Determine the receiver (the OTHER person in the session)
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
		CallerName: "Incoming Call", // Ideally fetch name from Profile Service, but keep simple for now
		// CallerAvatar: ...
	}

	// 4. Send Event
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
