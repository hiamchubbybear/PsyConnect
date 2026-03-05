package usecase

import (
	"consultationservice/internal/newsfeed/comment/domain"
	"consultationservice/internal/newsfeed/comment/repository"
	"context"
)


type PostRepository interface {
	UpdateEngagementCount(ctx context.Context, id, field string, delta int) error
	GetPostByID(ctx context.Context, id string) (interface{}, error)
}

type CreateCommentUseCase struct {
	commentRepo repository.CommentRepository
	postRepo    PostRepository
}

func NewCreateCommentUseCase(
	commentRepo repository.CommentRepository,
	postRepo PostRepository,
) *CreateCommentUseCase {
	return &CreateCommentUseCase{
		commentRepo: commentRepo,
		postRepo:    postRepo,
	}
}

func (uc *CreateCommentUseCase) Execute(ctx context.Context, postID, userID, content, parentCommentID string) (*domain.Comment, error) {
	comment := domain.NewComment(postID, userID, content, parentCommentID)

	err := uc.commentRepo.CreateComment(ctx, comment)
	if err != nil {
		return nil, err
	}

	
	if comment.IsRootComment() {
		go uc.postRepo.UpdateEngagementCount(context.Background(), postID, "comment_count", 1)
	}

	return comment, nil
}
