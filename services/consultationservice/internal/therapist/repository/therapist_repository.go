package repository

import (
	"consultationservice/internal/therapist/domain"
	"context"
)

type TherapistRepository interface {
	// Create creates a new therapist matching profile
	Create(ctx context.Context, therapist *domain.Therapist) error

	// GetByProfileID retrieves a therapist by profile ID
	GetByProfileID(ctx context.Context, profileID string) (*domain.Therapist, error)

	// GetAll retrieves all therapist matching profiles
	GetAll(ctx context.Context) ([]*domain.Therapist, error)

	// Update updates an existing therapist matching profile
	Update(ctx context.Context, profileID string, therapist *domain.Therapist) error

	// UpdateAvailability updates the availability status of a therapist
	UpdateAvailability(ctx context.Context, profileID string, isAvailable bool) error

	// Delete deletes a therapist matching profile
	Delete(ctx context.Context, profileID string) error

	// Search searches for therapists by name or specialization
	Search(ctx context.Context, query string, limit, skip int64) ([]*domain.Therapist, error)
}
