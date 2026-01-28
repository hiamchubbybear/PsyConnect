package domain

import "time"

const (
	SwipeStatusPending = "pending"
	SwipeStatusSwiped  = "swiped"
	SwipeStatusMatched = "matched"
)

type Swipe struct {
	ClientID    string    `bson:"client_id" json:"client_id"`
	TherapistID string    `bson:"therapist_id" json:"therapist_id"`
	Points      float32   `bson:"points" json:"points"`
	Reasons     []string  `bson:"reasons" json:"reasons"`
	Status      string    `bson:"status" json:"status"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
}

func NewSwipe(clientID, therapistID string, points float32, reasons []string) *Swipe {
	return &Swipe{
		ClientID:    clientID,
		TherapistID: therapistID,
		Points:      points,
		Reasons:     reasons,
		Status:      SwipeStatusPending,
		CreatedAt:   time.Now().UTC(),
	}
}
