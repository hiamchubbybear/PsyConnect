package handlers

import (
	"consultationservice/bootstrap"
	"consultationservice/internal/redis"
	"consultationservice/internal/repository"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type SocialHandler struct {
	env         *bootstrap.Env
	repoManager *repository.RepositoryManager
	redisClient redis.RedisStore
}

func NewSocialHandler(env *bootstrap.Env, repoManager *repository.RepositoryManager, redisClient redis.RedisStore) *SocialHandler {
	return &SocialHandler{
		env:         env,
		repoManager: repoManager,
		redisClient: redisClient,
	}
}

// Follow/Unfollow

func (h *SocialHandler) FollowUser(c *gin.Context) {
	followingID := c.Param("id")
	followerID := c.GetString("userID")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := h.repoManager.FollowRepo.FollowUser(ctx, followerID, followingID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Invalidate caches
	followingCacheKey := fmt.Sprintf("user:%s:following", followerID)
	followersCacheKey := fmt.Sprintf("user:%s:followers", followingID)
	h.redisClient.Delete(ctx, followingCacheKey)
	h.redisClient.Delete(ctx, followersCacheKey)

	c.JSON(http.StatusOK, gin.H{"message": "User followed successfully"})
}

func (h *SocialHandler) UnfollowUser(c *gin.Context) {
	followingID := c.Param("id")
	followerID := c.GetString("userID")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := h.repoManager.FollowRepo.UnfollowUser(ctx, followerID, followingID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Invalidate caches
	followingCacheKey := fmt.Sprintf("user:%s:following", followerID)
	followersCacheKey := fmt.Sprintf("user:%s:followers", followingID)
	h.redisClient.Delete(ctx, followingCacheKey)
	h.redisClient.Delete(ctx, followersCacheKey)

	c.JSON(http.StatusOK, gin.H{"message": "User unfollowed successfully"})
}

func (h *SocialHandler) GetFollowers(c *gin.Context) {
	userID := c.Param("id")
	limitStr := c.DefaultQuery("limit", "20")
	skipStr := c.DefaultQuery("skip", "0")

	limit, _ := strconv.ParseInt(limitStr, 10, 64)
	skip, _ := strconv.ParseInt(skipStr, 10, 64)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Try cache
	cacheKey := fmt.Sprintf("user:%s:followers:%d:%d", userID, limit, skip)
	var followers interface{}
	err := h.redisClient.Get(ctx, cacheKey, &followers)
	if err == nil {
		c.JSON(http.StatusOK, followers)
		return
	}

	// Get from DB
	followersList, err := h.repoManager.FollowRepo.GetFollowers(ctx, userID, limit, skip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get followers"})
		return
	}

	// Cache for 1 hour
	h.redisClient.Set(ctx, cacheKey, followersList)

	c.JSON(http.StatusOK, followersList)
}

func (h *SocialHandler) GetFollowing(c *gin.Context) {
	userID := c.Param("id")
	limitStr := c.DefaultQuery("limit", "20")
	skipStr := c.DefaultQuery("skip", "0")

	limit, _ := strconv.ParseInt(limitStr, 10, 64)
	skip, _ := strconv.ParseInt(skipStr, 10, 64)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Try cache
	cacheKey := fmt.Sprintf("user:%s:following:%d:%d", userID, limit, skip)
	var following interface{}
	err := h.redisClient.Get(ctx, cacheKey, &following)
	if err == nil {
		c.JSON(http.StatusOK, following)
		return
	}

	// Get from DB
	followingList, err := h.repoManager.FollowRepo.GetFollowing(ctx, userID, limit, skip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get following"})
		return
	}

	// Cache for 1 hour
	h.redisClient.Set(ctx, cacheKey, followingList)

	c.JSON(http.StatusOK, followingList)
}

// Bookmarks

func (h *SocialHandler) AddBookmark(c *gin.Context) {
	postID := c.Param("id")
	userID := c.GetString("userID")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := h.repoManager.BookmarkRepo.AddBookmark(ctx, userID, postID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("user:%s:bookmarks", userID)
	h.redisClient.Delete(ctx, cacheKey)

	c.JSON(http.StatusOK, gin.H{"message": "Post bookmarked"})
}

func (h *SocialHandler) RemoveBookmark(c *gin.Context) {
	postID := c.Param("id")
	userID := c.GetString("userID")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := h.repoManager.BookmarkRepo.RemoveBookmark(ctx, userID, postID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("user:%s:bookmarks", userID)
	h.redisClient.Delete(ctx, cacheKey)

	c.JSON(http.StatusOK, gin.H{"message": "Bookmark removed"})
}

func (h *SocialHandler) GetBookmarks(c *gin.Context) {
	userID := c.GetString("userID")
	limitStr := c.DefaultQuery("limit", "20")
	skipStr := c.DefaultQuery("skip", "0")

	limit, _ := strconv.ParseInt(limitStr, 10, 64)
	skip, _ := strconv.ParseInt(skipStr, 10, 64)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Try cache
	cacheKey := fmt.Sprintf("user:%s:bookmarks:%d:%d", userID, limit, skip)
	var bookmarks interface{}
	err := h.redisClient.Get(ctx, cacheKey, &bookmarks)
	if err == nil {
		c.JSON(http.StatusOK, bookmarks)
		return
	}

	// Get from DB
	bookmarksList, err := h.repoManager.BookmarkRepo.GetBookmarks(ctx, userID, limit, skip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get bookmarks"})
		return
	}

	// Cache for 30 minutes
	h.redisClient.Set(ctx, cacheKey, bookmarksList)

	c.JSON(http.StatusOK, bookmarksList)
}

// Share Post

func (h *SocialHandler) SharePost(c *gin.Context) {
	postID := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Increment share count
	err := h.repoManager.PostRepo.UpdateEngagementCount(ctx, postID, "share_count", 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to share post"})
		return
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("post:%s", postID)
	h.redisClient.Delete(ctx, cacheKey)

	c.JSON(http.StatusOK, gin.H{"message": "Post shared successfully"})
}
