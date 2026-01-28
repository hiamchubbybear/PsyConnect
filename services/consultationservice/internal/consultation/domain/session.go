package domain

import (
	"consultationservice/internal/payment/domain"
	"consultationservice/internal/utils"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ConsultationMode string

const (
	ConsultationModeOnline  ConsultationMode = "online"
	ConsultationModeOffline ConsultationMode = "offline"
)

type SessionStatus string

const (
	SessionStatusPending   SessionStatus = "pending"
	SessionStatusActive    SessionStatus = "active"
	SessionStatusCompleted SessionStatus = "completed"
	SessionStatusCancelled SessionStatus = "cancelled"
)

type CancelData struct {
	CancelledBy        *string `json:"cancelled_by,omitempty"`
	CancellationReason *string
}
type Session struct {
	// Info
	SessionID       string `json:"session_id" bson:"_id"`
	TherapistID     string `json:"therapist_id" bson:"therapist_id"`
	ClientID        string `json:"client_id" bson:"client_id"`
	SessionUserCode string `json:"session-user-code" bson:"session-user-code"`
	// Consultation detail
	Mode      ConsultationMode `json:"mode" bson:"mode"`
	StartTime time.Time        `json:	"start_time" bson:"start_time"`
	EndTime   time.Time        `json:"end_time" bson:"end_time"`
	Status    SessionStatus    `json:"status" bson:"status"`
	Timezone  string           `bson:"timezone" json:"timezone"`
	// Payment
	Price         float64              `json:"price" bson:"price"`
	PaymentStatus domain.PaymentStatus `json:"payment_status" bson:"payment_status"`
	PaymentID     *string              `json:"payment_id,omitempty" bson:"payment_id,omitempty"`
	CallSessionID *string              `json:"call_session_id,omitempty"`
	// Trace , log , notification
	CancelMetaData CancelData `json:"cancel_meta_data,omitempty"`
	ReminderSentAt time.Time  `json:"reminder_sent_at,omitempty"`
	ConversationID string     `json:"conversation_id" bson:"conversation_id"`
	CreatedAt      time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" bson:"updated_at"`
}

func NewSession(
	clientID, therapistID string,
	startTime, endTime time.Time,
	timezone string,
	price float64,
	mode ConsultationMode,
) (*Session, error) {

	now := time.Now().UTC()

	startUTC := startTime.UTC()
	endUTC := endTime.UTC()
	if _, err := time.LoadLocation(timezone); err != nil {
		return nil, fmt.Errorf("invalid timezone: %s", timezone)
	}

	s := &Session{
		SessionID:     uuid.New().String(),
		ClientID:      clientID,
		TherapistID:   therapistID,
		Mode:          mode,
		StartTime:     startUTC,
		EndTime:       endUTC,
		Status:        SessionStatusPending,
		PaymentStatus: domain.PaymentPending,
		PaymentID:     nil, // Nil for testing
		CallSessionID: nil, // Nil for testing
		Price:         price,
		Timezone:      timezone,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	s.SessionUserCode = s.GenerateConsultationCode()
	convID, err := s.GenerateConversationID()
	if err != nil {
		return nil, err
	}
	s.ConversationID = convID
	reminder := startUTC.Add(-30 * time.Minute)
	if reminder.After(now) {
		s.ReminderSentAt = reminder
	}
	return s, nil
}

func (s *Session) IsActive() bool {
	return s.Status == SessionStatusActive
}

func (s *Session) CanStartCall() bool {
	now := time.Now()
	return s.IsActive() &&
		now.After(s.StartTime.Add(-15*time.Minute)) &&
		now.Before(s.EndTime)
}

func (s *Session) Complete() {
	s.Status = SessionStatusCompleted
	s.UpdatedAt = time.Now().UTC()
}
func (s *Session) Activate() {
	s.Status = SessionStatusActive
	s.UpdatedAt = time.Now().UTC()
}
func (s *Session) Cancel() {
	s.Status = SessionStatusCancelled
	s.UpdatedAt = time.Now().UTC()
}
func (s *Session) GenerateConsultationCode() string {
	postFix := uuid.New().String()[:4]
	middleFix := s.StartTime.UTC().Format("20060102-1504")
	return fmt.Sprintf("CONS-%s-%s", middleFix, postFix)
}
func (s *Session) GenerateConversationID() (string, error) {
	util := utils.New()
	res, err := util.EncodeConversationId(s.ClientID, s.TherapistID)
	return res, err
}
