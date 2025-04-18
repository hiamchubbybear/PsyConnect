package model

type Therapist struct {
	ProfileId         string   `json:"profile_id,omitempty" bson:"profile_id"`
	Address           string   `json:"address,omitempty" bson:"address"`
	Languages         []string `json:"languages,omitempty" bson:"languages"`
	Specialization    []string `json:"specialization,omitempty" bson:"specialization"`
	ConsultationModes []string `json:"consultation_modes,omitempty" bson:"consultation_modes"`
	Experience        int      `json:"experience,omitempty" bson:"experience"`
	Rating            float64  `json:"rating,omitempty" bson:"rating"`
	PricePerSession   int      `json:"price_per_session,omitempty" bson:"price_per_session"`
	IsAvailable       bool     `json:"is_available,omitempty" bson:"is_available"`
	Availability      struct {
		Days      []string `json:"days,omitempty" bson:"days"`
		TimeSlots []string `json:"time_slots,omitempty" bson:"time_slots"`
	} `json:"availability,omitempty" bson:"availability"`
	CurrentSession []string `json:"current_session,omitempty" bson:"current_session"`
}
