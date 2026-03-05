package repository

import (
	"consultationservice/internal/model"
	"consultationservice/internal/newsfeed/post/domain"
	legacyRepo "consultationservice/internal/repository"
	"context"
)



type LegacyPostRepositoryAdapter struct {
	legacyRepo legacyRepo.PostRepository
}

func NewLegacyPostRepositoryAdapter(legacy legacyRepo.PostRepository) *LegacyPostRepositoryAdapter {
	return &LegacyPostRepositoryAdapter{
		legacyRepo: legacy,
	}
}

func (a *LegacyPostRepositoryAdapter) CreatePost(ctx context.Context, post *domain.Post) error {
	modelPost := convertDomainToModel(post)
	err := a.legacyRepo.CreatePost(ctx, modelPost)
	if err == nil {
		post.ID = modelPost.ID
	}
	return err
}

func (a *LegacyPostRepositoryAdapter) GetPostByID(ctx context.Context, id string) (*domain.Post, error) {
	modelPost, err := a.legacyRepo.GetPostByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return convertModelToDomain(modelPost), nil
}

func (a *LegacyPostRepositoryAdapter) UpdatePost(ctx context.Context, id string, post *domain.Post) error {
	return a.legacyRepo.UpdatePost(ctx, id, convertDomainToModel(post))
}

func (a *LegacyPostRepositoryAdapter) DeletePost(ctx context.Context, id string) error {
	return a.legacyRepo.DeletePost(ctx, id)
}

func (a *LegacyPostRepositoryAdapter) GetPostsByUser(ctx context.Context, userID string, limit, skip int64) ([]domain.Post, error) {
	modelPosts, err := a.legacyRepo.GetPostsByUser(ctx, userID, limit, skip)
	if err != nil {
		return nil, err
	}
	return convertModelSliceToDomain(modelPosts), nil
}

func (a *LegacyPostRepositoryAdapter) GetFeed(ctx context.Context, userIDs []string, excludeIDs []string, limit, skip int64) ([]domain.Post, error) {
	modelPosts, err := a.legacyRepo.GetFeed(ctx, userIDs, excludeIDs, limit, skip)
	if err != nil {
		return nil, err
	}
	return convertModelSliceToDomain(modelPosts), nil
}

func (a *LegacyPostRepositoryAdapter) SearchPosts(ctx context.Context, query string, limit, skip int64) ([]domain.Post, error) {
	modelPosts, err := a.legacyRepo.SearchPosts(ctx, query, limit, skip)
	if err != nil {
		return nil, err
	}
	return convertModelSliceToDomain(modelPosts), nil
}

func (a *LegacyPostRepositoryAdapter) GetTrendingPosts(ctx context.Context, limit int64) ([]domain.Post, error) {
	modelPosts, err := a.legacyRepo.GetTrendingPosts(ctx, limit)
	if err != nil {
		return nil, err
	}
	return convertModelSliceToDomain(modelPosts), nil
}

func (a *LegacyPostRepositoryAdapter) GetPostsByTag(ctx context.Context, tag string, limit, skip int64) ([]domain.Post, error) {
	modelPosts, err := a.legacyRepo.GetPostsByTag(ctx, tag, limit, skip)
	if err != nil {
		return nil, err
	}
	return convertModelSliceToDomain(modelPosts), nil
}

func (a *LegacyPostRepositoryAdapter) GetPostsByCategory(ctx context.Context, category string, limit, skip int64) ([]domain.Post, error) {
	modelPosts, err := a.legacyRepo.GetPostsByCategory(ctx, category, limit, skip)
	if err != nil {
		return nil, err
	}
	return convertModelSliceToDomain(modelPosts), nil
}

func (a *LegacyPostRepositoryAdapter) IncrementViewCount(ctx context.Context, id string) error {
	return a.legacyRepo.IncrementViewCount(ctx, id)
}

func (a *LegacyPostRepositoryAdapter) UpdateEngagementCount(ctx context.Context, id, field string, delta int) error {
	return a.legacyRepo.UpdateEngagementCount(ctx, id, field, delta)
}

func (a *LegacyPostRepositoryAdapter) GetPopularTags(ctx context.Context, limit int) ([]string, error) {
	return a.legacyRepo.GetPopularTags(ctx, limit)
}


func convertDomainToModel(d *domain.Post) *model.Post {
	return &model.Post{
		ID:            d.ID,
		Title:         d.Title,
		Content:       d.Content,
		AuthorID:      d.AuthorID,
		Tags:          d.Tags,
		Categories:    d.Categories,
		Media:         convertMediaDomainToModel(d.Media),
		Mentions:      d.Mentions,
		Hashtags:      d.Hashtags,
		Visibility:    d.Visibility,
		PostType:      d.PostType,
		ViewCount:     d.ViewCount,
		UpvoteCount:   d.UpvoteCount,
		DownvoteCount: d.DownvoteCount,
		CommentCount:  d.CommentCount,
		ShareCount:    d.ShareCount,
		UserVote:      d.UserVote,
		UserBookmark:  d.UserBookmark,
		IsDeleted:     d.IsDeleted,
		CreatedAt:     d.CreatedAt,
		UpdatedAt:     d.UpdatedAt,
	}
}

func convertModelToDomain(m *model.Post) *domain.Post {
	return &domain.Post{
		ID:            m.ID,
		Title:         m.Title,
		Content:       m.Content,
		AuthorID:      m.AuthorID,
		Tags:          m.Tags,
		Categories:    m.Categories,
		Media:         convertMediaModelToDomain(m.Media),
		Mentions:      m.Mentions,
		Hashtags:      m.Hashtags,
		Visibility:    m.Visibility,
		PostType:      m.PostType,
		ViewCount:     m.ViewCount,
		UpvoteCount:   m.UpvoteCount,
		DownvoteCount: m.DownvoteCount,
		CommentCount:  m.CommentCount,
		ShareCount:    m.ShareCount,
		UserVote:      m.UserVote,
		UserBookmark:  m.UserBookmark,
		IsDeleted:     m.IsDeleted,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

func convertModelSliceToDomain(models []model.Post) []domain.Post {
	result := make([]domain.Post, len(models))
	for i := range models {
		result[i] = *convertModelToDomain(&models[i])
	}
	return result
}

func convertMediaDomainToModel(media []domain.MediaAttachment) []model.MediaAttachment {
	result := make([]model.MediaAttachment, len(media))
	for i := range media {
		result[i] = model.MediaAttachment{
			Type:    media[i].Type,
			URL:     media[i].URL,
			Caption: media[i].Caption,
		}
	}
	return result
}

func convertMediaModelToDomain(media []model.MediaAttachment) []domain.MediaAttachment {
	result := make([]domain.MediaAttachment, len(media))
	for i := range media {
		result[i] = domain.MediaAttachment{
			Type:    media[i].Type,
			URL:     media[i].URL,
			Caption: media[i].Caption,
		}
	}
	return result
}


var _ PostRepository = (*LegacyPostRepositoryAdapter)(nil)
