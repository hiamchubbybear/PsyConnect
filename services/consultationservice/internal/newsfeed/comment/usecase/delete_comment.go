package usecase

import (
	"consultationservice/internal/newsfeed/comment/repository"
	"context"
)

type DeleteCommentUseCase struct {
	commentRepo repository.CommentRepository
	postRepo    PostRepository
}

func NewDeleteCommentUseCase(
	commentRepo repository.CommentRepository,
	postRepo PostRepository,
) *DeleteCommentUseCase {
	return &DeleteCommentUseCase{
		commentRepo: commentRepo,
		postRepo:    postRepo,
	}
}

func (uc *DeleteCommentUseCase) Execute(ctx context.Context, commentID, userID string) error {
	// Get comment to verify ownership
	comment, err := uc.commentRepo.GetCommentByID(ctx, commentID)
	if err != nil {
		return err
	}

	if comment.UserID != userID {
		return ErrNotAuthorized
	}

	err = uc.commentRepo.DeleteComment(ctx, commentID)
	if err != nil {
		return err
	}

	// Update post comment count (only for root comments)
	if comment.ParentCommentID == "" {
		go uc.postRepo.UpdateEngagementCount(context.Background(), comment.PostID, "comment_count", -1)
	}

	return nil
}
