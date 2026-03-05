package http

import (
	"consultationservice/internal/newsfeed/social/usecase"
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	followUserUC     *usecase.FollowUserUseCase
	unfollowUserUC   *usecase.UnfollowUserUseCase
	bookmarkPostUC   *usecase.BookmarkPostUseCase
	unbookmarkPostUC *usecase.UnbookmarkPostUseCase
}

func NewHandler(
	followUserUC *usecase.FollowUserUseCase,
	unfollowUserUC *usecase.UnfollowUserUseCase,
	bookmarkPostUC *usecase.BookmarkPostUseCase,
	unbookmarkPostUC *usecase.UnbookmarkPostUseCase,
) *Handler {
	return &Handler{
		followUserUC:     followUserUC,
		unfollowUserUC:   unfollowUserUC,
		bookmarkPostUC:   bookmarkPostUC,
		unbookmarkPostUC: unbookmarkPostUC,
	}
}

func (h *Handler) FollowUser(c *gin.Context) {
	followerID := c.GetString("userID")
	followingID := c.Param("userId")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := h.followUserUC.Execute(ctx, followerID, followingID)
	if err != nil {
		if errors.Is(err, usecase.ErrCannotFollowSelf) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot follow yourself"})
			return
		}
		if err.Error() == "already following" {
			c.JSON(http.StatusConflict, gin.H{"error": "Already following"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to follow user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Followed successfully"})
}

func (h *Handler) UnfollowUser(c *gin.Context) {
	followerID := c.GetString("userID")
	followingID := c.Param("userId")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := h.unfollowUserUC.Execute(ctx, followerID, followingID)
	if err != nil {
		if err.Error() == "not following" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not following"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unfollow user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Unfollowed successfully"})
}

func (h *Handler) BookmarkPost(c *gin.Context) {
	userID := c.GetString("userID")
	postID := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := h.bookmarkPostUC.Execute(ctx, userID, postID)
	if err != nil {
		if err.Error() == "already bookmarked" {
			c.JSON(http.StatusConflict, gin.H{"error": "Already bookmarked"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to bookmark post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bookmarked successfully"})
}

func (h *Handler) UnbookmarkPost(c *gin.Context) {
	userID := c.GetString("userID")
	postID := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := h.unbookmarkPostUC.Execute(ctx, userID, postID)
	if err != nil {
		if err.Error() == "bookmark not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Bookmark not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unbookmark post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Unbookmarked successfully"})
}


func (h *Handler) AddBookmark(c *gin.Context) {
	h.BookmarkPost(c)
}

func (h *Handler) RemoveBookmark(c *gin.Context) {
	h.UnbookmarkPost(c)
}

func (h *Handler) SharePost(c *gin.Context) {
	
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Share post not implemented"})
}

func (h *Handler) GetFollowers(c *gin.Context) {
	
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Get followers not implemented"})
}

func (h *Handler) GetFollowing(c *gin.Context) {
	
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Get following not implemented"})
}

func (h *Handler) GetBookmarks(c *gin.Context) {
	
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Get bookmarks not implemented"})
}
