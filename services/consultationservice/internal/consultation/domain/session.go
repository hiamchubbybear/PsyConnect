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
	
	SessionID       string `json:"session_id" bson:"_id"`
	TherapistID     string `json:"therapist_id" bson:"therapist_id"`
	ClientID        string `json:"client_id" bson:"client_id"`
	SessionUserCode string `json:"session-user-code" bson:"session-user-code"`
	
	Mode      ConsultationMode `json:"mode" bson:"mode"`
	StartTime time.Time        `json:"start_time" bson:"start_time"`
	EndTime   time.Time        `json:"end_time" bson:"end_time"`
	Status    SessionStatus    `json:"status" bson:"status"`
	
	Timezone      string        `bson:"timezone" json:"timezone"`
	ScheduledDate string        `bson:"scheduled_date" json:"scheduled_date"`
	LocationInfo  *LocationInfo `bson:"location_info,omitempty" json:"location_info,omitempty"`
	
	Price         float64              `json:"price" bson:"price"`
	PaymentStatus domain.PaymentStatus `json:"payment_status" bson:"payment_status"`
	PaymentID     *string              `json:"payment_id,omitempty" bson:"payment_id,omitempty"`
	RefundTraceID *string              `json:"refund_trace_id,omitempty" bson:"refund_trace_id,omitempty"`
	CallSessionID *string              `json:"call_session_id,omitempty"`
	
	CancelMetaData CancelData   `json:"cancel_meta_data,omitempty"`
	ReminderSentAt time.Time    `json:"reminder_sent_at,omitempty"`
	ConversationID string       `json:"conversation_id" bson:"conversation_id"`
	Logs           []SessionLog `json:"logs,omitempty" bson:"logs,omitempty"`
	CreatedAt      time.Time    `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at" bson:"updated_at"`
}

type LocationInfo struct {
	Link        string   `json:"link,omitempty" bson:"link,omitempty"`
	Passcode    string   `json:"passcode,omitempty" bson:"passcode,omitempty"`
	AddressLine string   `json:"address_line,omitempty" bson:"address_line,omitempty"`
	City        string   `json:"city,omitempty" bson:"city,omitempty"`
	Coordinates []string `json:"coordinates,omitempty" bson:"coordinates,omitempty"`
}

type SessionLog struct {
	Status    SessionStatus `json:"status" bson:"status"`
	Timestamp time.Time     `json:"timestamp" bson:"timestamp"`
	Message   string        `json:"message,omitempty" bson:"message,omitempty"`
}

func NewSession(
	clientID, therapistID string,
	startTime, endTime time.Time,
	timezone string,
	scheduledDate string,
	locationInfo *LocationInfo,
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
		PaymentID:     nil, 
		RefundTraceID: nil,
		CallSessionID: nil, 
		Price:         price,
		Timezone:      timezone,
		ScheduledDate: scheduledDate,
		LocationInfo:  locationInfo,
		Logs: []SessionLog{
			{
				Status:    SessionStatusPending,
				Timestamp: now,
				Message:   "Session initialized",
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
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

func (s *Session) AddLog(status SessionStatus, message string) {
	s.Logs = append(s.Logs, SessionLog{
		Status:    status,
		Timestamp: time.Now().UTC(),
		Message:   message,
	})
	s.UpdatedAt = time.Now().UTC()
}

func (s *Session) Complete() {
	s.Status = SessionStatusCompleted
	s.AddLog(SessionStatusCompleted, "Session marked as completed")
}

func (s *Session) Activate() {
	s.Status = SessionStatusActive
	s.AddLog(SessionStatusActive, "Session activated")
}

func (s *Session) Cancel(reason string) {
	s.Status = SessionStatusCancelled
	s.AddLog(SessionStatusCancelled, fmt.Sprintf("Session cancelled: %s", reason))
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
