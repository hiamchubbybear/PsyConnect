package model

type ClientSwipes struct {
	ClientId string           `json:"client_id,omitempty" bson:"client_id"`
	Swipes   []TherapistSwipe `json:"swipes,omitempty" bson:"swipes"`
}

type TherapistSwipe struct {
	TherapistId string   `json:"therapist_id,omitempty" bson:"therapist_id"`
	Points      float32  `json:"points,omitempty" bson:"points"`
	Reasons     []string `json:"reasons,omitempty" bson:"reasons"`
}
