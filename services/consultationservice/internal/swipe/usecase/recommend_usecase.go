package usecase

import (
	"context"

	clientDomain "consultationservice/internal/client/domain"
	clientRepository "consultationservice/internal/client/repository"

	"consultationservice/internal/dto"
	"consultationservice/internal/model"
	"consultationservice/internal/repository/infrastructure/external"

	swipeDomain "consultationservice/internal/swipe/domain"
	swipeRepository "consultationservice/internal/swipe/repository"

	therapistDomain "consultationservice/internal/therapist/domain"
	therapistRepository "consultationservice/internal/therapist/repository"
)

type RecommendUseCase struct {
	clientRepo    clientRepository.ClientRepository
	therapistRepo therapistRepository.TherapistRepository
	swipeRepo     swipeRepository.SwipeRepository
}

func NewRecommendUseCase(
	clientRepo clientRepository.ClientRepository,
	therapistRepo therapistRepository.TherapistRepository,
	swipeRepo swipeRepository.SwipeRepository,
) *RecommendUseCase {
	return &RecommendUseCase{
		clientRepo:    clientRepo,
		therapistRepo: therapistRepo,
		swipeRepo:     swipeRepo,
	}
}

func (uc *RecommendUseCase) TriggerUpdateV1(ctx context.Context, profileID string) error {
	client, err := uc.clientRepo.GetByProfileID(ctx, profileID)
	if err != nil {
		return err
	}

	therapists, err := uc.therapistRepo.GetAll(ctx)
	if err != nil {
		return err
	}

	rawClient := mapClientDomainToModel(client)
	var rawTherapists []model.TherapistV1
	for _, t := range therapists {
		rawTherapists = append(rawTherapists, mapTherapistDomainToModel(t))
	}

	swipes, err := external.RecommendationApiV1(dto.FilterRawDataV1{
		ClientRaw:    rawClient,
		TherapistRaw: rawTherapists,
	})
	if err != nil {
		return err
	}

	err = uc.swipeRepo.DeleteSwipesByClient(ctx, profileID)
	if err != nil {
		return err
	}

	var domainSwipes []*swipeDomain.Swipe
	for _, s := range swipes {
		domainSwipes = append(domainSwipes, &swipeDomain.Swipe{
			ClientID:    s.ClientId,
			TherapistID: s.TherapistId,
			Points:      float32(s.Points),
			Reasons:     s.Reasons,
			Status:      swipeDomain.SwipeStatusPending,
		})
	}

	return uc.swipeRepo.InsertSwipes(ctx, profileID, domainSwipes)
}

func (uc *RecommendUseCase) PopTop5V1(ctx context.Context, profileID string) ([]EnrichedSwipe, error) {
	swipes, err := uc.swipeRepo.GetTopSwipes(ctx, profileID, 5)
	if err != nil {
		return nil, err
	}

	var enrichedSwipes []EnrichedSwipe
	for _, s := range swipes {
		therapist, err := uc.therapistRepo.GetByProfileID(ctx, s.TherapistID)
		if err != nil {
			
			continue
		}

		enrichedSwipes = append(enrichedSwipes, EnrichedSwipe{
			TherapistV1: mapTherapistDomainToModel(therapist),
			Points:      s.Points,
			Reasons:     s.Reasons,
		})
	}

	return enrichedSwipes, nil
}

func mapClientDomainToModel(c *clientDomain.Client) model.Client {
	return model.Client{
		ProfileId:         c.ProfileID,
		Address:           c.Address,
		Languages:         c.Languages,
		IssueDetail:       c.IssueDetail,
		ConsultationModes: c.ConsultationModes,
		RangePrice:        c.RangePrice,
		Availability: struct {
			Days      []string `json:"days,omitempty" bson:"days" required:"true"`
			TimeSlots []string `json:"time_slots,omitempty" bson:"time_slots" required:"true"`
		}{
			Days:      c.Availability.Days,
			TimeSlots: c.Availability.TimeSlots,
		},
		PreferredTherapistGender:   c.PreferredTherapistGender,
		ExperienceLevel:            c.ExperienceLevel,
		TherapistSpecialization:    c.TherapistSpecialization,
		UrgencyLevel:               c.UrgencyLevel,
		SessionDuration:            c.SessionDuration,
		PreferredTherapistLanguage: c.PreferredTherapistLanguage,
		IsFlexibleWithSchedule:     c.IsFlexibleWithSchedule,
	}
}

func mapTherapistDomainToModel(t *therapistDomain.Therapist) model.TherapistV1 {
	degrees := make([]model.Degree, len(t.ProfessionalInfo.Degrees))
	for i, d := range t.ProfessionalInfo.Degrees {
		degrees[i] = model.Degree{
			Type:        d.Type,
			Field:       d.Field,
			Institution: d.Institution,
			Year:        d.Year,
		}
	}
	certs := make([]model.Certification, len(t.ProfessionalInfo.Certifications))
	for i, c := range t.ProfessionalInfo.Certifications {
		certs[i] = model.Certification{
			Name:   c.Name,
			Issuer: c.Issuer,
			Year:   c.Year,
		}
	}

	return model.TherapistV1{
		ProfileId:         t.ProfileID,
		Address:           t.Address,
		Languages:         t.Languages,
		Specialization:    t.Specialization,
		ConsultationModes: t.ConsultationModes,
		Experience:        t.Experience,
		Rating:            t.Rating,
		Currency:          t.Currency,
		RagePrice:         t.RagePrice,
		IsAvailable:       t.IsAvailable,
		Availability: struct {
			Days      []string `json:"days,omitempty" bson:"days"`
			TimeSlots []string `json:"time_slots,omitempty" bson:"time_slots"`
		}{
			Days:      t.Availability.Days,
			TimeSlots: t.Availability.TimeSlots,
		},
		CurrentSession: t.CurrentSession,
		MatchedClients: t.MatchedClients,
		AvatarOverride: t.AvatarOverride,
		Name:           t.Name,
		ProfessionalInfo: model.ProfessionalInfo{
			Title: model.ProfessionalTitle{
				Code:    t.ProfessionalInfo.Title.Code,
				Display: t.ProfessionalInfo.Title.Display,
			},
			Degrees:         degrees,
			Certifications:  certs,
			ExperienceYears: t.ProfessionalInfo.ExperienceYears,
		},
	}
}
