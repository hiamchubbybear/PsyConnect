package repository

import (
	"consultationservice/internal/consultation/domain"
	domain2 "consultationservice/internal/payment/domain"
	"context"
)


type SessionRepository interface {
	
	Create(ctx context.Context, session *domain.Session) error

	
	GetByID(ctx context.Context, id string) (*domain.Session, error)

	
	GetByTherapist(ctx context.Context, therapistID string) ([]*domain.Session, error)

	
	GetByClient(ctx context.Context, clientID string) ([]*domain.Session, error)

	
	GetAll(ctx context.Context) ([]*domain.Session, error)

	
	Update(ctx context.Context, session *domain.Session) error

	
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
	
	CheckTimeOverlap(ctx context.Context, sessions []*domain.Session, startTime, endTime string, day string) (bool, error)
}
