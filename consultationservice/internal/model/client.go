package model

type Client struct {
	ProfileId         string   `json:"profile_id,omitempty" bson:"profile_id" `
	Address           string   `json:"address,omitempty" bson:"address" required:"true"`
	Languages         []string `json:"languages,omitempty" bson:"languages" required:"true"`
	IssueDetail       []string `json:"issue_detail,omitempty" bson:"issue_detail" `
	ConsultationModes []string `json:"consultation_modes,omitempty" bson:"consultation_modes" required:"true"`
	RangePricePerHour int      `json:"range_price_per_hour,omitempty" bson:"range_price_per_hour" required:"true"`
	Availability      struct {
		Days      []string `json:"days,omitempty" bson:"days" required:"true"`
		TimeSlots []string `json:"time_slots,omitempty" bson:"time_slots" required:"true"`
	} `json:"availability,omitempty" bson:"availability" `
	PreferredTherapistGender   string   `json:"preferred_therapist_gender,omitempty" bson:"preferred_therapist_gender"`
	ExperienceLevel            string   `json:"experience_level,omitempty" bson:"experience_level"`
	TherapistSpecialization    []string `json:"therapist_specialization,omitempty" bson:"therapist_specialization"`
	UrgencyLevel               string   `json:"urgency_level,omitempty" bson:"urgency_level"`
	SessionDuration            int      `json:"session_duration,omitempty" bson:"session_duration"`
	PreferredTherapistLanguage []string `json:"preferred_therapist_language,omitempty" bson:"preferred_therapist_language"`
	IsFlexibleWithSchedule     bool     `json:"is_flexible_with_schedule,omitempty" bson:"is_flexible_with_schedule"`
	CurrentSession             []string `json:"current_session,omitempty" bson:"current_session"`
}
