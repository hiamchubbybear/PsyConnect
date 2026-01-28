package usecase

import (
	"consultationservice/internal/client/domain"
	"consultationservice/internal/client/repository"
	"context"
)

type CreateClientRequest struct {
	ProfileID                  string
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

type CreateClientUseCase struct {
	repo repository.ClientRepository
}

func NewCreateClientUseCase(repo repository.ClientRepository) *CreateClientUseCase {
	return &CreateClientUseCase{repo: repo}
}

func (uc *CreateClientUseCase) Execute(ctx context.Context, req CreateClientRequest) (*domain.Client, error) {
	availability := domain.Availability{
		Days:      req.AvailabilityDays,
		TimeSlots: req.AvailabilityTimeSlots,
	}

	client, err := domain.NewClient(
		req.ProfileID,
		req.Address,
		req.Languages,
		req.ConsultationModes,
		req.RangePrice,
		availability,
	)
	if err != nil {
		return nil, err
	}

	// Set optional fields
	client.IssueDetail = req.IssueDetail
	client.PreferredTherapistGender = req.PreferredTherapistGender
	client.ExperienceLevel = req.ExperienceLevel
	client.TherapistSpecialization = req.TherapistSpecialization
	client.UrgencyLevel = req.UrgencyLevel
	client.SessionDuration = req.SessionDuration
	client.PreferredTherapistLanguage = req.PreferredTherapistLanguage
	client.IsFlexibleWithSchedule = req.IsFlexibleWithSchedule

	if err := uc.repo.Create(ctx, client); err != nil {
		return nil, err
	}

	return client, nil
}
