package usecase

import (
	"consultationservice/internal/newsfeed/post/domain"
	"consultationservice/internal/newsfeed/post/repository"
	"context"
	"time"
)

// CreatePostUseCase handles post creation
type CreatePostUseCase struct {
	postRepo repository.PostRepository
}

func NewCreatePostUseCase(postRepo repository.PostRepository) *CreatePostUseCase {
	return &CreatePostUseCase{postRepo: postRepo}
}

func (uc *CreatePostUseCase) Execute(ctx context.Context, post *domain.Post) error {
	post.CreatedAt = time.Now().UTC()
	post.UpdatedAt = time.Now().UTC()
	return uc.postRepo.CreatePost(ctx, post)
}

// GetPostByIDUseCase retrieves a single post
type GetPostByIDUseCase struct {
	postRepo repository.PostRepository
}

func NewGetPostByIDUseCase(postRepo repository.PostRepository) *GetPostByIDUseCase {
	return &GetPostByIDUseCase{postRepo: postRepo}
}

func (uc *GetPostByIDUseCase) Execute(ctx context.Context, id string) (*domain.Post, error) {
	return uc.postRepo.GetPostByID(ctx, id)
}

// UpdatePostUseCase handles post updates with ownership verification
type UpdatePostUseCase struct {
	postRepo repository.PostRepository
}

func NewUpdatePostUseCase(postRepo repository.PostRepository) *UpdatePostUseCase {
	return &UpdatePostUseCase{postRepo: postRepo}
}

func (uc *UpdatePostUseCase) Execute(ctx context.Context, id, userID string, updates *domain.Post) error {
	// Verify ownership
	existing, err := uc.postRepo.GetPostByID(ctx, id)
	if err != nil {
		return err
	}
	if existing.AuthorID != userID {
		return ErrNotAuthorized
	}

	updates.UpdatedAt = time.Now().UTC()
	return uc.postRepo.UpdatePost(ctx, id, updates)
}

// DeletePostUseCase handles post deletion with ownership verification
type DeletePostUseCase struct {
	postRepo repository.PostRepository
}

func NewDeletePostUseCase(postRepo repository.PostRepository) *DeletePostUseCase {
	return &DeletePostUseCase{postRepo: postRepo}
}

func (uc *DeletePostUseCase) Execute(ctx context.Context, id, userID string) error {
	// Verify ownership
	existing, err := uc.postRepo.GetPostByID(ctx, id)
	if err != nil {
		return err
	}
	if existing.AuthorID != userID {
		return ErrNotAuthorized
	}

	return uc.postRepo.DeletePost(ctx, id)
}

// GetFeedUseCase retrieves user feed
type GetFeedUseCase struct {
	postRepo repository.PostRepository
}

func NewGetFeedUseCase(postRepo repository.PostRepository) *GetFeedUseCase {
	return &GetFeedUseCase{postRepo: postRepo}
}

func (uc *GetFeedUseCase) Execute(ctx context.Context, userIDs []string, excludeIDs []string, limit, skip int64) ([]domain.Post, error) {
	return uc.postRepo.GetFeed(ctx, userIDs, excludeIDs, limit, skip)
}

// GetTrendingPostsUseCase retrieves trending posts
type GetTrendingPostsUseCase struct {
	postRepo repository.PostRepository
}

func NewGetTrendingPostsUseCase(postRepo repository.PostRepository) *GetTrendingPostsUseCase {
	return &GetTrendingPostsUseCase{postRepo: postRepo}
}

func (uc *GetTrendingPostsUseCase) Execute(ctx context.Context, limit int64) ([]domain.Post, error) {
	return uc.postRepo.GetTrendingPosts(ctx, limit)
}

// SearchPostsUseCase handles post search
type SearchPostsUseCase struct {
	postRepo repository.PostRepository
}

func NewSearchPostsUseCase(postRepo repository.PostRepository) *SearchPostsUseCase {
	return &SearchPostsUseCase{postRepo: postRepo}
}

func (uc *SearchPostsUseCase) Execute(ctx context.Context, query string, limit, skip int64) ([]domain.Post, error) {
	return uc.postRepo.SearchPosts(ctx, query, limit, skip)
}

// GetUserPostsUseCase retrieves posts by a specific user
type GetUserPostsUseCase struct {
	postRepo repository.PostRepository
}

func NewGetUserPostsUseCase(postRepo repository.PostRepository) *GetUserPostsUseCase {
	return &GetUserPostsUseCase{postRepo: postRepo}
}

func (uc *GetUserPostsUseCase) Execute(ctx context.Context, userID string, limit, skip int64) ([]domain.Post, error) {
	return uc.postRepo.GetPostsByUser(ctx, userID, limit, skip)
}

// GetPostsByTagUseCase retrieves posts by tag
type GetPostsByTagUseCase struct {
	postRepo repository.PostRepository
}

func NewGetPostsByTagUseCase(postRepo repository.PostRepository) *GetPostsByTagUseCase {
	return &GetPostsByTagUseCase{postRepo: postRepo}
}

func (uc *GetPostsByTagUseCase) Execute(ctx context.Context, tag string, limit, skip int64) ([]domain.Post, error) {
	return uc.postRepo.GetPostsByTag(ctx, tag, limit, skip)
}

// GetPostsByCategoryUseCase retrieves posts by category
type GetPostsByCategoryUseCase struct {
	postRepo repository.PostRepository
}

func NewGetPostsByCategoryUseCase(postRepo repository.PostRepository) *GetPostsByCategoryUseCase {
	return &GetPostsByCategoryUseCase{postRepo: postRepo}
}

func (uc *GetPostsByCategoryUseCase) Execute(ctx context.Context, category string, limit, skip int64) ([]domain.Post, error) {
	return uc.postRepo.GetPostsByCategory(ctx, category, limit, skip)
}

// IncrementViewCountUseCase handles view count increment
type IncrementViewCountUseCase struct {
	postRepo repository.PostRepository
}

func NewIncrementViewCountUseCase(postRepo repository.PostRepository) *IncrementViewCountUseCase {
	return &IncrementViewCountUseCase{postRepo: postRepo}
}

func (uc *IncrementViewCountUseCase) Execute(ctx context.Context, id string) error {
	return uc.postRepo.IncrementViewCount(ctx, id)
}

// GetPopularTagsUseCase retrieves most used hashtags
type GetPopularTagsUseCase struct {
	postRepo repository.PostRepository
}

func NewGetPopularTagsUseCase(postRepo repository.PostRepository) *GetPopularTagsUseCase {
	return &GetPopularTagsUseCase{postRepo: postRepo}
}

func (uc *GetPopularTagsUseCase) Execute(ctx context.Context, limit int) ([]string, error) {
	return uc.postRepo.GetPopularTags(ctx, limit)
}

// Errors
var ErrNotAuthorized = &NotAuthorizedError{}

type NotAuthorizedError struct{}

func (e *NotAuthorizedError) Error() string {
	return "not authorized to perform this action"
}
