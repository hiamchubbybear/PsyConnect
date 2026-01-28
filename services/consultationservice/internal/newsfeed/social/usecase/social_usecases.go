package usecase

import (
	"consultationservice/internal/newsfeed/social/repository"
	"context"
)

type FollowUserUseCase struct {
	followRepo repository.FollowRepository
}

func NewFollowUserUseCase(followRepo repository.FollowRepository) *FollowUserUseCase {
	return &FollowUserUseCase{
		followRepo: followRepo,
	}
}

func (uc *FollowUserUseCase) Execute(ctx context.Context, followerID, followingID string) error {
	if followerID == followingID {
		return ErrCannotFollowSelf
	}
	return uc.followRepo.FollowUser(ctx, followerID, followingID)
}

type UnfollowUserUseCase struct {
	followRepo repository.FollowRepository
}

func NewUnfollowUserUseCase(followRepo repository.FollowRepository) *UnfollowUserUseCase {
	return &UnfollowUserUseCase{
		followRepo: followRepo,
	}
}

func (uc *UnfollowUserUseCase) Execute(ctx context.Context, followerID, followingID string) error {
	return uc.followRepo.UnfollowUser(ctx, followerID, followingID)
}

type BookmarkPostUseCase struct {
	bookmarkRepo repository.BookmarkRepository
}

func NewBookmarkPostUseCase(bookmarkRepo repository.BookmarkRepository) *BookmarkPostUseCase {
	return &BookmarkPostUseCase{
		bookmarkRepo: bookmarkRepo,
	}
}

func (uc *BookmarkPostUseCase) Execute(ctx context.Context, userID, postID string) error {
	return uc.bookmarkRepo.CreateBookmark(ctx, userID, postID)
}

type UnbookmarkPostUseCase struct {
	bookmarkRepo repository.BookmarkRepository
}

func NewUnbookmarkPostUseCase(bookmarkRepo repository.BookmarkRepository) *UnbookmarkPostUseCase {
	return &UnbookmarkPostUseCase{
		bookmarkRepo: bookmarkRepo,
	}
}

func (uc *UnbookmarkPostUseCase) Execute(ctx context.Context, userID, postID string) error {
	return uc.bookmarkRepo.DeleteBookmark(ctx, userID, postID)
}

var ErrCannotFollowSelf = &CannotFollowSelfError{}

type CannotFollowSelfError struct{}

func (e *CannotFollowSelfError) Error() string {
	return "cannot follow yourself"
}
