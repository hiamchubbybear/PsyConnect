package repository

import (
	"consultationservice/internal/therapist/domain"
	"context"
)

type TherapistRepository interface {
	
	Create(ctx context.Context, therapist *domain.Therapist) error

	
	GetByProfileID(ctx context.Context, profileID string) (*domain.Therapist, error)

	
	GetAll(ctx context.Context) ([]*domain.Therapist, error)

	
	Update(ctx context.Context, profileID string, therapist *domain.Therapist) error

	
	UpdateAvailability(ctx context.Context, profileID string, isAvailable bool) error

	
	Delete(ctx context.Context, profileID string) error

	
	Search(ctx context.Context, query string, limit, skip int64) ([]*domain.Therapist, error)
}
