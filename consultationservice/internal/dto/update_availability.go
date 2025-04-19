package dto

type Availability struct {
	Days      []string `json:"days,omitempty" bson:"days"`
	TimeSlots []string `json:"time_slots,omitempty" bson:"time_slots"`
}
