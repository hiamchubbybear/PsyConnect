package dto

type TherapistUpdateRequest struct {
	Address           string   `json:"address,omitempty" bson:"address"`
	Languages         []string `json:"languages,omitempty" bson:"languages"`
	Specialization    []string `json:"specialization,omitempty" bson:"specialization"`
	ConsultationModes []string `json:"consultation_modes,omitempty" bson:"consultation_modes"`
	Experience        int      `json:"experience,omitempty" bson:"experience"`
	Rating            float64  `json:"rating,omitempty" bson:"rating"`
	PricePerSession   int      `json:"price_per_session,omitempty" bson:"price_per_session"`
}
