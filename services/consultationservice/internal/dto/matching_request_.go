package dto

import "time"

type MatchingRequest struct {
	TherapistId string    `json:"therapistId" `
	Message     string    `json:"message"`
	RequestedAt time.Time `json:"requested_at,omitempty"`
}
