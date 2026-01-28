package usecase

import (
	"consultationservice/internal/newsfeed/comment/domain"
	"consultationservice/internal/newsfeed/comment/repository"
	"context"
)

type GetCommentsUseCase struct {
	commentRepo repository.CommentRepository
}

func NewGetCommentsUseCase(commentRepo repository.CommentRepository) *GetCommentsUseCase {
	return &GetCommentsUseCase{
		commentRepo: commentRepo,
	}
}

func (uc *GetCommentsUseCase) Execute(ctx context.Context, postID string, maxDepth int, limit, skip int64) ([]domain.Comment, error) {
	comments, err := uc.commentRepo.GetCommentsTree(ctx, postID, maxDepth, limit, skip)
	if err != nil {
		return nil, err
	}

	// Ensure comments is never nil
	if comments == nil {
		comments = []domain.Comment{}
	}

	return comments, nil
}
