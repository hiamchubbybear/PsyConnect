package model

import "time"

// Deprecated: Replace ClientSwipeV1
type ClientSwipe struct {
	ClientId    string    `json:"client_id" bson:"client_id"`
	TherapistId string    `json:"therapist_id" bson:"therapist_id"`
	Points      float32   `json:"points" bson:"points"`
	Reasons     []string  `json:"reasons" bson:"reasons"`
	Status      string    `json:"status" bson:"status"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
}
type ClientSwipeV1 struct {
	ClientId    string    `json:"client_id" bson:"client_id"`
	TherapistId string    `json:"therapist_id" bson:"therapist_id"`
	Points      float32   `json:"points" bson:"points"`
	Reasons     []string  `json:"reasons" bson:"reasons"`
	Status      string    `json:"status" bson:"status"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
}
