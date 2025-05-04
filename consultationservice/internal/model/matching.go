package model

type MatchingModel struct {
	Address                    string   `json:"address" bson:"address"`
	Languages                  []string `json:"languages" bson:"languages"`
	IssueDetail                []string `json:"issue_detail" bson:"issue_detail"`
	ConsultationModes          []string `json:"consultation_modes" bson:"consultation_modes"`
	RangePricePerHour          int      `json:"range_price_per_hour" bson:"price_per_session"`
	PreferredTherapistGender   string   `json:"preferred_therapist_gender" bson:"gender"`
	ExperienceLevel            string   `json:"experience_level" bson:"experience_level"`
	TherapistSpecialization    []string `json:"therapist_specialization" bson:"specialization"`
	UrgencyLevel               string   `json:"urgency_level" bson:"urgency_level"`
	SessionDuration            int      `json:"session_duration" bson:"session_duration"`
	PreferredTherapistLanguage []string `json:"preferred_therapist_language" bson:"languages"`
	IsFlexibleWithSchedule     bool     `json:"is_flexible_with_schedule" bson:"is_flexible_with_schedule"`
	Availability               struct {
		Days      []string `json:"days" bson:"availability.days"`
		TimeSlots []string `json:"time_slots" bson:"availability.time_slots"`
	} `json:"availability" bson:"availability"`
}
