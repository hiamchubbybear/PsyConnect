package usecase

import (
	"consultationservice/internal/repository"
	"context"
)

// PostRepositoryAdapter adapts legacy PostRepository to usecase interface
type PostRepositoryAdapter struct {
	postRepo repository.PostRepository
}

func NewPostRepositoryAdapter(postRepo repository.PostRepository) *PostRepositoryAdapter {
	return &PostRepositoryAdapter{
		postRepo: postRepo,
	}
}

func (a *PostRepositoryAdapter) UpdateEngagementCount(ctx context.Context, id, field string, delta int) error {
	return a.postRepo.UpdateEngagementCount(ctx, id, field, delta)
}

func (a *PostRepositoryAdapter) GetPostByID(ctx context.Context, id string) (interface{}, error) {
	return a.postRepo.GetPostByID(ctx, id)
}

// Ensure adapter implements the interface
var _ PostRepository = (*PostRepositoryAdapter)(nil)
