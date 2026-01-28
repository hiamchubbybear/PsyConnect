package domain

import (
	"errors"
	"time"
)

type Availability struct {
	Days      []string
	TimeSlots []string
}

type Client struct {
	ProfileID                  string
	Address                    string
	Languages                  []string
	IssueDetail                []string
	ConsultationModes          []string
	RangePrice                 int
	Availability               Availability
	PreferredTherapistGender   string
	ExperienceLevel            string
	TherapistSpecialization    []string
	UrgencyLevel               string
	SessionDuration            int
	PreferredTherapistLanguage []string
	IsFlexibleWithSchedule     bool
	CreatedAt                  time.Time
	UpdatedAt                  time.Time
}

func NewClient(
	profileID string,
	address string,
	languages []string,
	consultationModes []string,
	rangePrice int,
	availability Availability,
) (*Client, error) {
	if profileID == "" {
		return nil, errors.New("profile ID is required")
	}
	if address == "" {
		return nil, errors.New("address is required")
	}
	if len(languages) == 0 {
		return nil, errors.New("at least one language is required")
	}
	if len(consultationModes) == 0 {
		return nil, errors.New("at least one consultation mode is required")
	}
	if rangePrice <= 0 {
		return nil, errors.New("range price must be greater than 0")
	}

	now := time.Now()
	return &Client{
		ProfileID:         profileID,
		Address:           address,
		Languages:         languages,
		ConsultationModes: consultationModes,
		RangePrice:        rangePrice,
		Availability:      availability,
		CreatedAt:         now,
		UpdatedAt:         now,
	}, nil
}

func (c *Client) Update() {
	c.UpdatedAt = time.Now()
}
