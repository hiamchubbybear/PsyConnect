package dto

import (
	"consultationservice/internal/consultation/domain"
	"time"
)

type SessionRequest struct {
	TherapistID string                  `json:"therapist_id"`
	ClientID    string                  `json:"client_id"`
	Mode        domain.ConsultationMode `json:"mode"`
	StartTime   time.Time               `json:"start_time"`
	EndTime     time.Time               `json:"end_time"`
	Timezone    string                  `json:"timezone"`
	Price       float64                 `json:"price"`
}
