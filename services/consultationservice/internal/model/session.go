package model

import (
	"consultationservice/internal/enum"
	"time"

	"github.com/google/uuid"
)

type LocationInfo struct {
	Link        string   `json:"link,omitempty" bson:"link,omitempty"`
	Passcode    string   `json:"passcode,omitempty" bson:"passcode,omitempty"`
	AddressLine string   `json:"address_line,omitempty" bson:"address_line,omitempty"`
	City        string   `json:"city,omitempty" bson:"city,omitempty"`
	Coordinates []string `json:"coordinates,omitempty" bson:"coordinates,omitempty"` // [longitude, latitude]
}

type TimeInfo struct {
	Timezone      string `json:"timezone" bson:"timezone"`
	ScheduledDate string `json:"scheduled_date" bson:"scheduled_date"` // YYYY-MM-DD
}

type SessionLog struct {
	Status    string    `json:"status" bson:"status"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
	Message   string    `json:"message,omitempty" bson:"message,omitempty"`
}

type Session struct {
	SessionId     string                `json:"session_id" bson:"_id"`
	TherapistID   string                `json:"therapist_id" bson:"therapist_id"`
	ClientID      string                `json:"client_id" bson:"client_id"`
	Mode          enum.ConsultationMode `json:"mode" bson:"mode"`
	TimeInfo      TimeInfo              `json:"time_info" bson:"time_info"`
	LocationInfo  *LocationInfo         `json:"location_info,omitempty" bson:"location_info,omitempty"`
	StartTime     time.Time             `json:"start_time" bson:"start_time"`
	EndTime       time.Time             `json:"end_time" bson:"end_time"`
	Status        string                `json:"status" bson:"status"`
	Price         float64               `json:"price" bson:"price"`
	TransactionID string                `json:"transaction_id,omitempty" bson:"transaction_id,omitempty"`
	Gateway       string                `json:"gateway,omitempty" bson:"gateway,omitempty"`
	PaymentStatus string                `json:"payment_status,omitempty" bson:"payment_status,omitempty"`
	RefundTraceID string                `json:"refund_trace_id,omitempty" bson:"refund_trace_id,omitempty"`
	Logs          []SessionLog          `json:"logs,omitempty" bson:"logs,omitempty"`
	CreatedAt     time.Time             `json:"created_at" bson:"created_at"`
}

func NewSession(clientID, therapistID string, price float64, startTime, endTime time.Time, mode enum.ConsultationMode, timeInfo TimeInfo, locationInfo *LocationInfo) *Session {
	now := time.Now().UTC()
	return &Session{
		SessionId:     uuid.New().String(),
		ClientID:      clientID,
		TherapistID:   therapistID,
		Status:        MatchStatusActive, // Consider changing to PAYMENT_PENDING when fully integrating payments
		PaymentStatus: "PENDING",
		StartTime:     startTime,
		EndTime:       endTime,
		Mode:          mode,
		TimeInfo:      timeInfo,
		LocationInfo:  locationInfo,
		Price:         price,
		Logs: []SessionLog{
			{
				Status:    "CREATED",
				Timestamp: now,
				Message:   "Session initialized",
			},
		},
		CreatedAt: now,
	}
}
