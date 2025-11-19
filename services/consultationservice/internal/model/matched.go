package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	MatchStatusActive    = "active"
	MatchStatusClosed    = "closed"
	MatchStatusCancelled = "cancelled"
)

type Match struct {
	MatchID           string     `bson:"match_id" json:"match_id"`
	ClientID          string     `bson:"client_id" json:"client_id"`
	TherapistID       string     `bson:"therapist_id" json:"therapist_id"`
	MatchedAt         time.Time  `bson:"matched_at" json:"matched_at"`
	Status            string     `bson:"status" json:"status"`
	Source            string     `bson:"source" json:"source"`
	SwipeScore        float64    `bson:"swipe_score" json:"swipe_score"`
	Reasons           []string   `bson:"reasons" json:"reasons"`
	ClientFeedback    string     `bson:"client_feedback,omitempty" json:"client_feedback,omitempty"`
	TherapistFeedback string     `bson:"therapist_feedback,omitempty" json:"therapist_feedback,omitempty"`
	ClosedAt          *time.Time `bson:"closed_at,omitempty" json:"closed_at,omitempty"`
	UpdatedAt         time.Time  `bson:"updated_at" json:"updated_at"`
}

func NewMatch(clientID, therapistID, source string, score float64, reasons []string) *Match {
	return &Match{
		MatchID:     uuid.New().String(),
		ClientID:    clientID,
		TherapistID: therapistID,
		MatchedAt:   time.Now().UTC(),
		Status:      MatchStatusActive,
		Source:      source,
		SwipeScore:  score,
		Reasons:     reasons,
		UpdatedAt:   time.Now().UTC(),
	}
}
