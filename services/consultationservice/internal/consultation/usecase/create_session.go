package usecase

import (
	"consultationservice/internal/consultation/domain"
	"consultationservice/internal/consultation/repository"
	"consultationservice/internal/kafka"
	"context"
	"errors"
	"time"
)

type CreateSessionRequest struct {
	TherapistID   string
	ClientID      string
	Mode          domain.ConsultationMode
	StartTime     time.Time
	EndTime       time.Time
	Price         float64
	TimeZone      string
	ScheduledDate string
	LocationInfo  *domain.LocationInfo
}

type CreateSessionUseCase struct {
	sessionRepo repository.SessionRepository
	producer    *kafka.Producer
}

func NewCreateSessionUseCase(sessionRepo repository.SessionRepository, producer *kafka.Producer) *CreateSessionUseCase {
	return &CreateSessionUseCase{
		sessionRepo: sessionRepo,
		producer:    producer,
	}
}

func (uc *CreateSessionUseCase) Execute(ctx context.Context, req CreateSessionRequest) (*domain.Session, error) {

	if req.TherapistID == "" {
		return nil, errors.New("therapist ID is required")
	}
	if req.ClientID == "" {
		return nil, errors.New("client ID is required")
	}
	if req.StartTime.IsZero() || req.EndTime.IsZero() {
		return nil, errors.New("start time and end time are required")
	}
	if req.EndTime.Before(req.StartTime) {
		return nil, errors.New("end time must be after start time")
	}
	if req.Price <= 0 {
		return nil, errors.New("price must be greater than 0")
	}
	if req.TimeZone == "" {
		return nil, errors.New("time zone can not null")
	}
	therapistSessions, err := uc.sessionRepo.GetByTherapist(ctx, req.TherapistID)
	if err != nil {
		return nil, err
	}

	clientSessions, err := uc.sessionRepo.GetByClient(ctx, req.ClientID)
	if err != nil {
		return nil, err
	}

	allSessions := append(therapistSessions, clientSessions...)

	startTimeStr := req.StartTime.Format("15:04")
	endTimeStr := req.EndTime.Format("15:04")
	day := req.StartTime.Weekday().String()

	overlap, err := uc.sessionRepo.CheckTimeOverlap(ctx, allSessions, startTimeStr, endTimeStr, day)
	if err != nil {
		return nil, err
	}
	if overlap {
		return nil, errors.New("session time overlaps with existing sessions")
	}
	if req.ScheduledDate == "" {
		req.ScheduledDate = req.StartTime.Format("2006-01-02")
	}

	// Payment nil first
	session, err := domain.NewSession(
		req.ClientID,
		req.TherapistID,
		req.StartTime,
		req.EndTime,
		req.TimeZone,
		req.ScheduledDate,
		req.LocationInfo,
		req.Price,
		req.Mode,
	)

	if err := uc.sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}

	// 5. Send Event
	if uc.producer != nil {
		_ = uc.producer.SendSessionEvent("notification.push.consultation-created", session)
	}

	return session, nil
}
