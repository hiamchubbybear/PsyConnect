package dto

import "consultationservice/internal/enum"

type SessionRequest struct {
	TherapistID       string                  `json:"therapist_id" bson:"therapist_id"`
	ClientID          string                  `json:"client_id" bson:"client_id"`
	Languages         []string                `json:"languages" bson:"languages"`
	Specialization    []string                `json:"specialization" bson:"specialization"`
	ConsultationModes []enum.ConsultationMode `json:"consultation_modes" bson:"consultation_modes"`
	PricePerHour      float64                 `json:"price_per_hour" bson:"price_per_hour"`
	SessionTime       SessionTime             `json:"session_time" bson:"session_time"`
}

type SessionTime struct {
	Day       enum.DayOfWeek `json:"day" bson:"day"`
	StartTime string         `json:"start_time" bson:"start_time"`
	EndTime   string         `json:"end_time" bson:"end_time"`
}

type DeleteSessionRequest struct {
	TherapistID string `json:"therapist_id" bson:"therapist_id"`
	ClientID    string `json:"client_id" bson:"client_id"`
	SessionID   string `json:"session_id" bson:"_id"`
}
