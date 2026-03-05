package usecase

import (
	"consultationservice/internal/therapist/domain"
	"consultationservice/internal/therapist/repository"
	"context"
)

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

type CreateTherapistRequest struct {
	ProfileID         string
	Address           string
	Languages         []string
	Specialization    []string
	ConsultationModes []string
	RagePrice         int
	Currency          string

	
	AvailabilityDays      []string
	AvailabilityTimeSlots []string

	
	Experience       int
	Rating           float64
	AvatarOverride   string
	Name             string
	ProfessionalInfo ProfessionalInfo
}

type CreateTherapistUseCase struct {
	repo repository.TherapistRepository
}

func NewCreateTherapistUseCase(repo repository.TherapistRepository) *CreateTherapistUseCase {
	return &CreateTherapistUseCase{repo: repo}
}

func (uc *CreateTherapistUseCase) Execute(ctx context.Context, req CreateTherapistRequest) (*domain.Therapist, error) {
	therapist, err := domain.NewTherapist(
		req.ProfileID,
		req.Address,
		req.Languages,
		req.Specialization,
		req.ConsultationModes,
		req.RagePrice,
		req.Currency,
	)
	if err != nil {
		return nil, err
	}

	therapist.Availability = domain.Availability{
		Days:      req.AvailabilityDays,
		TimeSlots: req.AvailabilityTimeSlots,
	}
	therapist.Experience = req.Experience
	therapist.Rating = req.Rating
	therapist.AvatarOverride = req.AvatarOverride
	therapist.Name = req.Name

	
	therapist.ProfessionalInfo = domain.ProfessionalInfo{
		Title: domain.ProfessionalTitle{
			Code:    req.ProfessionalInfo.Title.Code,
			Display: req.ProfessionalInfo.Title.Display,
		},
		ExperienceYears: req.ProfessionalInfo.ExperienceYears,
	}

	for _, d := range req.ProfessionalInfo.Degrees {
		therapist.ProfessionalInfo.Degrees = append(therapist.ProfessionalInfo.Degrees, domain.Degree{
			Type:        d.Type,
			Field:       d.Field,
			Institution: d.Institution,
			Year:        d.Year,
		})
	}

	for _, c := range req.ProfessionalInfo.Certifications {
		therapist.ProfessionalInfo.Certifications = append(therapist.ProfessionalInfo.Certifications, domain.Certification{
			Name:   c.Name,
			Issuer: c.Issuer,
			Year:   c.Year,
		})
	}

	if err := uc.repo.Create(ctx, therapist); err != nil {
		return nil, err
	}

	return therapist, nil
}
