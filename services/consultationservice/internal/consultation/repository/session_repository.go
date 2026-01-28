package repository

import (
	"consultationservice/internal/consultation/domain"
	domain2 "consultationservice/internal/payment/domain"
	"context"
)

// SessionRepository defines the interface for session data access
type SessionRepository interface {
	// Create creates a new session
	Create(ctx context.Context, session *domain.Session) error

	// GetByID retrieves a session by ID
	GetByID(ctx context.Context, id string) (*domain.Session, error)

	// GetByTherapist retrieves all sessions for a therapist
	GetByTherapist(ctx context.Context, therapistID string) ([]*domain.Session, error)

	// GetByClient retrieves all sessions for a client
	GetByClient(ctx context.Context, clientID string) ([]*domain.Session, error)

	// GetAll retrieves all sessions
	GetAll(ctx context.Context) ([]*domain.Session, error)

	// Update updates an existing session
	Update(ctx context.Context, session *domain.Session) error

	// Delete deletes a session by ID
	Delete(ctx context.Context, id string) error
	UpdateStatus(
		ctx context.Context,
		sessionID string,
		status domain.SessionStatus,
	) error

	UpdatePayment(
		ctx context.Context,
		sessionID string,
		paymentStatus domain2.PaymentStatus,
		paymentID *string,
	) error

	AttachCallSession(
		ctx context.Context,
		sessionID string,
		callSessionID string,
	) error

	Cancel(
		ctx context.Context,
		sessionID string,
		cancelMeta domain.CancelData,
	) error
	// CheckTimeOverlap checks if a new session overlaps with existing sessions
	CheckTimeOverlap(ctx context.Context, sessions []*domain.Session, startTime, endTime string, day string) (bool, error)
}
