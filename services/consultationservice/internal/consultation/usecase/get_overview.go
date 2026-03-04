package usecase

import (
	"consultationservice/internal/consultation/domain"
	"consultationservice/internal/consultation/repository"
	"context"
	"time"
)

type OverviewData struct {
	UpcomingCount       int               `json:"upcoming_count"`
	CompletedCount      int               `json:"completed_count"`
	CancelledCount      int               `json:"cancelled_count"`
	PendingPaymentCount int               `json:"pending_payment_count"`
	TotalSpent          float64           `json:"total_spent"`
	NextSession         *domain.Session   `json:"next_session"`
	RecentSessions      []*domain.Session `json:"recent_sessions"`
}

type GetOverviewUseCase struct {
	sessionRepo repository.SessionRepository
}

func NewGetOverviewUseCase(sessionRepo repository.SessionRepository) *GetOverviewUseCase {
	return &GetOverviewUseCase{
		sessionRepo: sessionRepo,
	}
}

func (uc *GetOverviewUseCase) Execute(ctx context.Context, profileID string, role string) (*OverviewData, error) {
	var sessions []*domain.Session
	var err error

	if role == "role.therapist" {
		sessions, err = uc.sessionRepo.GetByTherapist(ctx, profileID)
	} else {
		sessions, err = uc.sessionRepo.GetByClient(ctx, profileID)
	}

	if err != nil {
		return nil, err
	}

	now := time.Now()
	var upcomingCount, completedCount, cancelledCount, pendingPaymentCount int
	var totalSpent float64
	var nextSession *domain.Session

	for _, s := range sessions {
		switch s.Status {
		case domain.SessionStatusCompleted:
			completedCount++
			totalSpent += s.Price
		case domain.SessionStatusCancelled:
			cancelledCount++
		default:
			if s.StartTime.After(now) {
				upcomingCount++
				if nextSession == nil || s.StartTime.Before(nextSession.StartTime) {
					nextSession = s
				}
			}
		}

		if s.PaymentStatus == "unpaid" || s.PaymentStatus == "" {
			pendingPaymentCount++
		}
	}

	return &OverviewData{
		UpcomingCount:       upcomingCount,
		CompletedCount:      completedCount,
		CancelledCount:      cancelledCount,
		PendingPaymentCount: pendingPaymentCount,
		TotalSpent:          totalSpent,
		NextSession:         nextSession,
		RecentSessions:      sessions,
	}, nil
}
