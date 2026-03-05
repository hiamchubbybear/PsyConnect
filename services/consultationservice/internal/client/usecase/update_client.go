package usecase

import (
	"consultationservice/internal/client/domain"
	"consultationservice/internal/client/repository"
	"context"
)

type UpdateClientRequest struct {
	Address                    string
	Languages                  []string
	IssueDetail                []string
	ConsultationModes          []string
	RangePrice                 int
	AvailabilityDays           []string
	AvailabilityTimeSlots      []string
	PreferredTherapistGender   string
	ExperienceLevel            string
	TherapistSpecialization    []string
	UrgencyLevel               string
	SessionDuration            int
	PreferredTherapistLanguage []string
	IsFlexibleWithSchedule     bool
}

type UpdateClientUseCase struct {
	repo repository.ClientRepository
}

func NewUpdateClientUseCase(repo repository.ClientRepository) *UpdateClientUseCase {
	return &UpdateClientUseCase{repo: repo}
}

func (uc *UpdateClientUseCase) Execute(ctx context.Context, profileID string, req UpdateClientRequest) error {
	
	client, err := uc.repo.GetByProfileID(ctx, profileID)
	if err != nil {
		return err
	}

	
	if req.Address != "" {
		client.Address = req.Address
	}
	if len(req.Languages) > 0 {
		client.Languages = req.Languages
	}
	if len(req.ConsultationModes) > 0 {
		client.ConsultationModes = req.ConsultationModes
	}
	if req.RangePrice > 0 {
		client.RangePrice = req.RangePrice
	}
	if len(req.AvailabilityDays) > 0 || len(req.AvailabilityTimeSlots) > 0 {
		client.Availability = domain.Availability{
			Days:      req.AvailabilityDays,
			TimeSlots: req.AvailabilityTimeSlots,
		}
	}

	
	client.IssueDetail = req.IssueDetail
	client.PreferredTherapistGender = req.PreferredTherapistGender
	client.ExperienceLevel = req.ExperienceLevel
	client.TherapistSpecialization = req.TherapistSpecialization
	client.UrgencyLevel = req.UrgencyLevel
	client.SessionDuration = req.SessionDuration
	client.PreferredTherapistLanguage = req.PreferredTherapistLanguage
	client.IsFlexibleWithSchedule = req.IsFlexibleWithSchedule

	return uc.repo.Update(ctx, profileID, client)
}
