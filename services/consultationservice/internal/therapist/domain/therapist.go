package domain

import (
	"errors"
	"time"
)

type Availability struct {
	Days      []string
	TimeSlots []string
}

type ProfessionalTitle struct {
	Code    string
	Display string
}

type Degree struct {
	Type        string
	Field       string
	Institution string
	Year        int
}

type Certification struct {
	Name   string
	Issuer string
	Year   int
}

type ProfessionalInfo struct {
	Title           ProfessionalTitle
	Degrees         []Degree
	Certifications  []Certification
	ExperienceYears int
}

type Therapist struct {
	ProfileID         string
	Address           string
	Languages         []string
	Specialization    []string
	ConsultationModes []string
	Experience        int
	Rating            float64
	Currency          string
	RagePrice         int
	IsAvailable       bool
	Availability      Availability
	CurrentSession    []string
	MatchedClients    []string
	AvatarOverride    string
	Name              string
	ProfessionalInfo  ProfessionalInfo
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func NewTherapist(
	profileID string,
	address string,
	languages []string,
	specialization []string,
	consultationModes []string,
	ragePrice int,
	currency string,
) (*Therapist, error) {
	if profileID == "" {
		return nil, errors.New("profile ID is required")
	}
	if address == "" {
		return nil, errors.New("address is required")
	}
	if len(languages) == 0 {
		return nil, errors.New("at least one language is required")
	}
	if len(specialization) == 0 {
		return nil, errors.New("at least one specialization is required")
	}
	if ragePrice <= 0 {
		return nil, errors.New("price must be greater than 0")
	}

	now := time.Now()
	return &Therapist{
		ProfileID:         profileID,
		Address:           address,
		Languages:         languages,
		Specialization:    specialization,
		ConsultationModes: consultationModes,
		RagePrice:         ragePrice,
		Currency:          currency,
		IsAvailable:       true,
		MatchedClients:    []string{},
		CurrentSession:    []string{},
		CreatedAt:         now,
		UpdatedAt:         now,
	}, nil
}

func (t *Therapist) Update() {
	t.UpdatedAt = time.Now()
}

func (t *Therapist) SetAvailability(isAvailable bool) {
	t.IsAvailable = isAvailable
	t.UpdatedAt = time.Now()
}
