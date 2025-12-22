package handlers

import (
	"consultationservice/bootstrap"
	"consultationservice/internal/model"
	"consultationservice/internal/redis"
	"consultationservice/internal/repository"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

type CommentHandler struct {
	env         *bootstrap.Env
	repoManager *repository.RepositoryManager
	redisClient redis.RedisStore
}

func NewCommentHandler(
	env *bootstrap.Env,
	repoManager *repository.RepositoryManager,
	redisClient redis.RedisStore,
) *CommentHandler {
	return &CommentHandler{
		env:         env,
		repoManager: repoManager,
		redisClient: redisClient,
	}
}

func (h *CommentHandler) CreateComment(c *gin.Context) {
	postID := c.Param("id")
	userID := c.GetString("userID")

	// DEBUG: Log incoming request
	fmt.Printf("=== CREATE COMMENT DEBUG ===\n")
	fmt.Printf("PostID: %s\n", postID)
	fmt.Printf("UserID: %s\n", userID)
	fmt.Printf("Headers: %+v\n", c.Request.Header)
	fmt.Printf("============================\n")

	var req struct {
		Content         string `json:"content" binding:"required"`
		ParentCommentID string `json:"parent_comment_id,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("ERROR: Binding failed: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content is required"})
		return
	}

	fmt.Printf("Request Content: %s\n", req.Content)
	fmt.Printf("Request ParentCommentID: %s\n", req.ParentCommentID)

	// Manual validation for content length (Gin doesn't support max tag for strings)
	if len(req.Content) < 1 || len(req.Content) > 500 {
		fmt.Printf("ERROR: Content length invalid: %d characters\n", len(req.Content))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content must be 1-500 characters"})
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
		// Check for max depth error
		if err.Error() == "invalid index value" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot reply more than 3 levels deep"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	// Update post comment count (only for root comments)
	if req.ParentCommentID == "" {
		err = h.repoManager.PostRepo.UpdateEngagementCount(ctx, postID, "comment_count", 1)
		if err != nil {
			fmt.Printf("Failed to update comment count: %v\n", err)
		}
	}

	// Invalidate cache
	// Current Redis wrapper doesn't support wildcards, so we use a simpler strategy
	// We'll use a specific key for the most common view, and for others we'll rely on TTL
	cacheKey := fmt.Sprintf("post:%s", postID)
	h.redisClient.Delete(ctx, cacheKey)

	// Invalidate the comment tree cache for common depths
	for _, d := range []int{0, 1, 2, 5} {
		treeCacheKey := fmt.Sprintf("post:%s:comments:tree:%d:20:0", postID, d)
		h.redisClient.Delete(ctx, treeCacheKey)
	}

	fmt.Printf("DEBUG: Invalidated caches for post: %s\n", postID)

	// Emit Notification Event
	go func() {
		post, err := h.repoManager.PostRepo.GetPostByID(context.Background(), postID)
		if err != nil || post == nil {
			return
		}

		// Don't notify self
		if post.AuthorID == userID {
			return
		}

		commenterProfile, err := h.repoManager.GrpcProfile.GetProfile(userID)
		commenterName := "Someone"
		if err == nil && commenterProfile != nil {
			commenterName = fmt.Sprintf("%s %s", commenterProfile.FirstName, commenterProfile.LastName)
		}

		notificationData := map[string]interface{}{
			"userId":        post.AuthorID,
			"commenterName": commenterName,
			"postId":        postID,
			"commentId":     comment.ID.Hex(),
		}
		jsonData, _ := json.Marshal(notificationData)
		h.repoManager.Kafka.SendToTopic("notification.social.post-comment", string(jsonData))
	}()

	c.JSON(http.StatusCreated, comment)
}

func (h *CommentHandler) GetComments(c *gin.Context) {
	postID := c.Param("id")
	limitStr := c.DefaultQuery("limit", "20")
	skipStr := c.DefaultQuery("skip", "0")
	depthStr := c.DefaultQuery("depth", "2")

	limit, _ := strconv.ParseInt(limitStr, 10, 64)
	skip, _ := strconv.ParseInt(skipStr, 10, 64)
	maxDepth, _ := strconv.Atoi(depthStr)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Try cache first
	cacheKey := fmt.Sprintf("post:%s:comments:tree:%d:%d:%d", postID, maxDepth, limit, skip)
	var cachedComments []model.Comment
	err := h.redisClient.Get(ctx, cacheKey, &cachedComments)
	if err == nil {
		log.Printf("[CommentAPI] CACHE HIT for key: %s (Found %d comments)", cacheKey, len(cachedComments))
		c.JSON(http.StatusOK, gin.H{
			"comments": cachedComments,
			"total":    len(cachedComments),
			"source":   "cache",
		})
		return
	}

	log.Printf("[CommentAPI] CACHE MISS for key: %s. Fetching from DB...", cacheKey)

	// Get comment tree from DB
	comments, err := h.repoManager.CommentRepo.GetCommentsTree(ctx, postID, maxDepth, limit, skip)
	if err != nil {
		log.Printf("[CommentAPI] ERROR: GetCommentsTree failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comments"})
		return
	}

	log.Printf("[CommentAPI] DB SUCCESS: Found %d comments for post %s (maxDepth: %d)", len(comments), postID, maxDepth)

	// Ensure comments is never nil
	if comments == nil {
		comments = []model.Comment{}
	}

	// Populate author info
	comments = h.populateCommentAuthors(ctx, comments)

	// Cache result (5 minutes)
	h.redisClient.Set(ctx, cacheKey, comments)

	c.JSON(http.StatusOK, gin.H{
		"comments": comments,
		"total":    len(comments),
		"source":   "db",
	})
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

// fetchUserProfile fetches user profile from profile service with Redis caching
func (h *CommentHandler) fetchUserProfile(ctx context.Context, userID string) *model.AuthorInfo {
	// Try cache first
	cacheKey := fmt.Sprintf("user:profile:%s", userID)
	var cached model.AuthorInfo
	if err := h.redisClient.Get(ctx, cacheKey, &cached); err == nil {
		return &cached
	}

	// Fetch from profile service via gRPC
	profile, err := h.repoManager.GrpcProfile.GetProfile(userID)
	if err != nil {
		log.Printf("Failed to fetch profile for %s via gRPC: %v", userID, err)
		return nil
	}

	if profile == nil {
		log.Printf("Profile not found for user %s", userID)
		return nil
	}

	authorInfo := &model.AuthorInfo{
		ID:        profile.ProfileId,
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		Avatar:    profile.Avatar,
	}

	// Cache for 5 minutes
	h.redisClient.Set(ctx, cacheKey, authorInfo)

	return authorInfo
}

// populateCommentAuthors recursively populates author info for comments and replies
func (h *CommentHandler) populateCommentAuthors(ctx context.Context, comments []model.Comment) []model.Comment {
	for i := range comments {
		// Set AuthorID for frontend compatibility
		comments[i].AuthorID = comments[i].UserID

		// Fetch and populate author info
		if comments[i].UserID != "" {
			comments[i].Author = h.fetchUserProfile(ctx, comments[i].UserID)
		}

		// Recursively populate replies
		if len(comments[i].Replies) > 0 {
			comments[i].Replies = h.populateCommentAuthors(ctx, comments[i].Replies)
		}
	}
	return comments
}
