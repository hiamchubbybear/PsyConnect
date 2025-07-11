package dto

type TherapistCreateRequest struct {
	ProfileId         string   `json:"profile_id,omitempty" bson:"profile_id"`
	Address           string   `json:"address,omitempty" bson:"address"`
	Languages         []string `json:"languages,omitempty" bson:"languages"`
	Specialization    []string `json:"specialization,omitempty" bson:"specialization"`
	ConsultationModes []string `json:"consultation_modes,omitempty" bson:"consultation_modes"`
	Experience        int      `json:"experience,omitempty" bson:"experience"`
	RagePrice         int      `json:"rage_price,omitempty" bson:"rage-price"`
	IsAvailable       bool     `json:"is_available,omitempty" bson:"is_available"`
	Currency          string   `json:"currency" bson:"currency"`
	Availability      struct {
		Days      []string `json:"days,omitempty" bson:"days"`
		TimeSlots []string `json:"time_slots,omitempty" bson:"time_slots"`
	} `json:"availability,omitempty" bson:"availability"`
	CurrentSession []string `json:"current_session,omitempty" bson:"current_session"`
	MatchedClients []string `json:"matched_clients,omitempty" bson:"matched_clients"`
}
