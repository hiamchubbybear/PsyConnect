package usecase

import (
	"consultationservice/internal/therapist/domain"
	"consultationservice/internal/therapist/repository"
	"context"
)

type UpdateTherapistRequest struct {
	Address               string
	Languages             []string
	Specialization        []string
	ConsultationModes     []string
	RagePrice             int
	Currency              string
	AvailabilityDays      []string
	AvailabilityTimeSlots []string
	Experience            int
	AvatarOverride        string
	Name                  string
	ProfessionalInfo      *ProfessionalInfo 
}

type UpdateTherapistUseCase struct {
	repo repository.TherapistRepository
}

func NewUpdateTherapistUseCase(repo repository.TherapistRepository) *UpdateTherapistUseCase {
	return &UpdateTherapistUseCase{repo: repo}
}

func (uc *UpdateTherapistUseCase) Execute(ctx context.Context, profileID string, req UpdateTherapistRequest) error {
	therapist, err := uc.repo.GetByProfileID(ctx, profileID)
	if err != nil {
		return err
	}

	if req.Address != "" {
		therapist.Address = req.Address
	}
	if len(req.Languages) > 0 {
		therapist.Languages = req.Languages
	}
	if len(req.Specialization) > 0 {
		therapist.Specialization = req.Specialization
	}
	if len(req.ConsultationModes) > 0 {
		therapist.ConsultationModes = req.ConsultationModes
	}
	if req.RagePrice > 0 {
		therapist.RagePrice = req.RagePrice
	}
	if req.Currency != "" {
		therapist.Currency = req.Currency
	}
	if len(req.AvailabilityDays) > 0 || len(req.AvailabilityTimeSlots) > 0 {
		therapist.Availability.Days = req.AvailabilityDays
		therapist.Availability.TimeSlots = req.AvailabilityTimeSlots
	}
	if req.Experience > 0 {
		therapist.Experience = req.Experience
	}
	if req.AvatarOverride != "" {
		therapist.AvatarOverride = req.AvatarOverride
	}
	if req.Name != "" {
		therapist.Name = req.Name
	}

	if req.ProfessionalInfo != nil {
		therapist.ProfessionalInfo = domain.ProfessionalInfo{
			Title: domain.ProfessionalTitle{
				Code:    req.ProfessionalInfo.Title.Code,
				Display: req.ProfessionalInfo.Title.Display,
			},
			ExperienceYears: req.ProfessionalInfo.ExperienceYears,
		}

		therapist.ProfessionalInfo.Degrees = make([]domain.Degree, 0)
		for _, d := range req.ProfessionalInfo.Degrees {
			therapist.ProfessionalInfo.Degrees = append(therapist.ProfessionalInfo.Degrees, domain.Degree{
				Type:        d.Type,
				Field:       d.Field,
				Institution: d.Institution,
				Year:        d.Year,
			})
		}

		therapist.ProfessionalInfo.Certifications = make([]domain.Certification, 0)
		for _, c := range req.ProfessionalInfo.Certifications {
			therapist.ProfessionalInfo.Certifications = append(therapist.ProfessionalInfo.Certifications, domain.Certification{
				Name:   c.Name,
				Issuer: c.Issuer,
				Year:   c.Year,
			})
		}
	}

	return uc.repo.Update(ctx, profileID, therapist)
}

func (uc *UpdateTherapistUseCase) ExecuteAvailability(ctx context.Context, profileID string, isAvailable bool) error {
	return uc.repo.UpdateAvailability(ctx, profileID, isAvailable)
}
