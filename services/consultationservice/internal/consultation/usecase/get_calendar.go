package usecase

import (
	"consultationservice/internal/consultation/domain"
	"consultationservice/internal/consultation/repository"
	"context"
	"time"
)

type CalendarEvent struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Start       string `json:"start"`
	End         string `json:"end"`
	Mode        string `json:"mode"`
	Status      string `json:"status"`
	Color       string `json:"color"`
	TherapistID string `json:"therapist_id,omitempty"`
	ClientID    string `json:"client_id,omitempty"`
}

type GetCalendarUseCase struct {
	sessionRepo repository.SessionRepository
}

func NewGetCalendarUseCase(sessionRepo repository.SessionRepository) *GetCalendarUseCase {
	return &GetCalendarUseCase{
		sessionRepo: sessionRepo,
	}
}

func (uc *GetCalendarUseCase) Execute(ctx context.Context, profileID string, role string, fromStr, toStr string) ([]CalendarEvent, error) {
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

	var events []CalendarEvent
	modeColors := map[string]string{
		"online":    "#667eea",
		"in_person": "#48bb78",
		"phone":     "#ed8936",
		"chat":      "#4299e1",
	}

	for _, s := range sessions {
		startStr := s.StartTime.Format(time.RFC3339)

		// Apply date range filter
		if fromStr != "" && startStr < fromStr {
			continue
		}
		if toStr != "" && startStr > toStr {
			continue
		}

		color := modeColors[string(s.Mode)]
		if color == "" {
			color = "#667eea"
		}

		events = append(events, CalendarEvent{
			ID:          s.SessionID,
			Title:       "Consultation Session",
			Start:       startStr,
			End:         s.EndTime.Format(time.RFC3339),
			Mode:        string(s.Mode),
			Status:      string(s.Status),
			Color:       color,
			TherapistID: s.TherapistID,
			ClientID:    s.ClientID,
		})
	}

	if events == nil {
		events = []CalendarEvent{}
	}

	return events, nil
}
