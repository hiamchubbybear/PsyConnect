package handlers

import (
	"consultationservice/bootstrap"
	"consultationservice/internal/model"
	"consultationservice/internal/redis"
	"consultationservice/internal/repository"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	env         *bootstrap.Env
	repoManager *repository.RepositoryManager
	redisClient redis.RedisStore
}

func NewCommentHandler(env *bootstrap.Env, repoManager *repository.RepositoryManager, redisClient redis.RedisStore) *CommentHandler {
	return &CommentHandler{
		env:         env,
		repoManager: repoManager,
		redisClient: redisClient,
	}
}

func (h *CommentHandler) CreateComment(c *gin.Context) {
	postID := c.Param("id")
	userID := c.GetString("userID")

	var req struct {
		Content         string `json:"content" binding:"required"`
		ParentCommentID string `json:"parent_comment_id,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	comment := &model.Comment{
		PostID:          postID,
		UserID:          userID,
		Content:         req.Content,
		ParentCommentID: req.ParentCommentID,
	}

	err := h.repoManager.CommentRepo.CreateComment(ctx, comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	// Update post comment count (only for top-level comments)
	if req.ParentCommentID == "" {
		err = h.repoManager.PostRepo.UpdateEngagementCount(ctx, postID, "comment_count", 1)
		if err != nil {
			fmt.Printf("Failed to update comment count: %v\n", err)
		}
	}

	// Invalidate caches
	cacheKey := fmt.Sprintf("post:%s", postID)
	h.redisClient.Delete(ctx, cacheKey)

	commentsCacheKey := fmt.Sprintf("post:%s:comments", postID)
	h.redisClient.Delete(ctx, commentsCacheKey)

	c.JSON(http.StatusCreated, comment)
}

func (h *CommentHandler) GetComments(c *gin.Context) {
	postID := c.Param("id")
	limitStr := c.DefaultQuery("limit", "20")
	skipStr := c.DefaultQuery("skip", "0")

	limit, _ := strconv.ParseInt(limitStr, 10, 64)
	skip, _ := strconv.ParseInt(skipStr, 10, 64)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Try cache first
	cacheKey := fmt.Sprintf("post:%s:comments:%d:%d", postID, limit, skip)
	var comments []model.Comment
	err := h.redisClient.Get(ctx, cacheKey, &comments)
	if err == nil {
		c.JSON(http.StatusOK, comments)
		return
	}

	// Get from DB
	comments, err = h.repoManager.CommentRepo.GetCommentsByPost(ctx, postID, limit, skip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comments"})
		return
	}

	// Cache result
	h.redisClient.Set(ctx, cacheKey, comments)

	c.JSON(http.StatusOK, comments)
}

func (h *CommentHandler) UpdateComment(c *gin.Context) {
	commentID := c.Param("id")
	userID := c.GetString("userID")

	var req struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get comment to verify ownership
	comment, err := h.repoManager.CommentRepo.GetCommentByID(ctx, commentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	if comment.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	err = h.repoManager.CommentRepo.UpdateComment(ctx, commentID, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
		return
	}

	// Invalidate cache
	commentsCacheKey := fmt.Sprintf("post:%s:comments", comment.PostID)
	h.redisClient.Delete(ctx, commentsCacheKey)

	c.JSON(http.StatusOK, gin.H{"message": "Comment updated"})
}

func (h *CommentHandler) DeleteComment(c *gin.Context) {
	commentID := c.Param("id")
	userID := c.GetString("userID")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get comment to verify ownership
	comment, err := h.repoManager.CommentRepo.GetCommentByID(ctx, commentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	if comment.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	err = h.repoManager.CommentRepo.DeleteComment(ctx, commentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	// Update post comment count
	if comment.ParentCommentID == "" {
		err = h.repoManager.PostRepo.UpdateEngagementCount(ctx, comment.PostID, "comment_count", -1)
		if err != nil {
			fmt.Printf("Failed to update comment count: %v\n", err)
		}
	}

	// Invalidate cache
	postCacheKey := fmt.Sprintf("post:%s", comment.PostID)
	h.redisClient.Delete(ctx, postCacheKey)

	commentsCacheKey := fmt.Sprintf("post:%s:comments", comment.PostID)
	h.redisClient.Delete(ctx, commentsCacheKey)

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted"})
}

func (h *CommentHandler) GetReplies(c *gin.Context) {
	parentCommentID := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	replies, err := h.repoManager.CommentRepo.GetReplies(ctx, parentCommentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get replies"})
		return
	}

	c.JSON(http.StatusOK, replies)
}
