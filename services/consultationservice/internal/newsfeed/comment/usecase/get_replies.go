package usecase

import (
	"consultationservice/internal/newsfeed/comment/repository"
	"context"
)

type GetRepliesUseCase struct {
	commentRepo repository.CommentRepository
}

func NewGetRepliesUseCase(commentRepo repository.CommentRepository) *GetRepliesUseCase {
	return &GetRepliesUseCase{
		commentRepo: commentRepo,
	}
}

func (uc *GetRepliesUseCase) Execute(ctx context.Context, parentCommentID string) (interface{}, error) {
	return uc.commentRepo.GetReplies(ctx, parentCommentID)
}
