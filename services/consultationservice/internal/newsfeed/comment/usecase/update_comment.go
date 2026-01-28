package usecase

import (
	"consultationservice/internal/newsfeed/comment/repository"
	"context"
)

type UpdateCommentUseCase struct {
	commentRepo repository.CommentRepository
}

func NewUpdateCommentUseCase(commentRepo repository.CommentRepository) *UpdateCommentUseCase {
	return &UpdateCommentUseCase{
		commentRepo: commentRepo,
	}
}

func (uc *UpdateCommentUseCase) Execute(ctx context.Context, commentID, userID, content string) error {
	// Get comment to verify ownership
	comment, err := uc.commentRepo.GetCommentByID(ctx, commentID)
	if err != nil {
		return err
	}

	if comment.UserID != userID {
		return ErrNotAuthorized
	}

	return uc.commentRepo.UpdateComment(ctx, commentID, content)
}

var ErrNotAuthorized = &NotAuthorizedError{}

type NotAuthorizedError struct{}

func (e *NotAuthorizedError) Error() string {
	return "not authorized"
}
