package repository

import (
	"consultationservice/internal/client/domain"
	"context"
)

type ClientRepository interface {
	// Create creates a new client matching profile
	Create(ctx context.Context, client *domain.Client) error

	// GetByProfileID retrieves a client by profile ID
	GetByProfileID(ctx context.Context, profileID string) (*domain.Client, error)

	// GetAll retrieves all client matching profiles
	GetAll(ctx context.Context) ([]*domain.Client, error)

	// Update updates an existing client matching profile
	Update(ctx context.Context, profileID string, client *domain.Client) error

	// Delete deletes a client matching profile
	Delete(ctx context.Context, profileID string) error
}
