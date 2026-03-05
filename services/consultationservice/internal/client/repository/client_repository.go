package repository

import (
	"consultationservice/internal/client/domain"
	"context"
)

type ClientRepository interface {
	
	Create(ctx context.Context, client *domain.Client) error

	
	GetByProfileID(ctx context.Context, profileID string) (*domain.Client, error)

	
	GetAll(ctx context.Context) ([]*domain.Client, error)

	
	Update(ctx context.Context, profileID string, client *domain.Client) error

	
	Delete(ctx context.Context, profileID string) error
}
