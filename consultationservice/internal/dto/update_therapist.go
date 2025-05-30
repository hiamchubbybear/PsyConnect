package dto

type TherapistUpdateRequest struct {
	Address           string   `json:"address,omitempty" bson:"address"`
	Languages         []string `json:"languages,omitempty" bson:"languages"`
	Specialization    []string `json:"specialization,omitempty" bson:"specialization"`
	ConsultationModes []string `json:"consultation_modes,omitempty" bson:"consultation_modes"`
	Experience        int      `json:"experience,omitempty" bson:"experience"`
	Rating            float64  `json:"rating,omitempty" bson:"rating"`
	RagePrice         int      `json:"rage_price,omitempty" bson:"rage-price"`
	IsAvailable       bool     `json:"is_available,omitempty" bson:"is_available"`
	Availability      struct {
		Days      []string `json:"days,omitempty" bson:"days"`
		TimeSlots []string `json:"time_slots,omitempty" bson:"time_slots"`
	} `json:"availability,omitempty" bson:"availability"`
}
