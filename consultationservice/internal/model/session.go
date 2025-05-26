package model

import (
	"consultationservice/internal/enum"
	"time"
)

type Session struct {
	SessionId   string                `json:"session_id" bson:"_id"`
	TherapistID string                `json:"therapist_id" bson:"therapist_id"`
	ClientID    string                `json:"client_id" bson:"client_id"`
	Mode        enum.ConsultationMode `json:"mode" bson:"mode"`
	StartTime   time.Time             `json:"start_time" bson:"start_time"`
	EndTime     time.Time             `json:"end_time" bson:"end_time"`
	Status      string                `json:"status" bson:"status"`
	Price       float64               `json:"price" bson:"price"`
	CreatedAt   time.Time             `json:"created_at" bson:"created_at"`
}
