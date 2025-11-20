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

type PostHandler struct {
	env         *bootstrap.Env
	repoManager *repository.RepositoryManager
	redisClient redis.RedisStore
}

func NewPostHandler(env *bootstrap.Env, repoManager *repository.RepositoryManager, redisClient redis.RedisStore) *PostHandler {
	return &PostHandler{
		env:         env,
		repoManager: repoManager,
		redisClient: redisClient,
	}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	var post model.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		fmt.Printf("ERROR: Failed to bind JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Debug: Print received post data
	fmt.Printf("DEBUG: Received post data: %+v\n", post)

	// Get author ID from auth middleware
	authorID := c.GetString("userID")
	fmt.Printf("DEBUG: Author ID from auth: %s\n", authorID)

	post.AuthorID = authorID
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Printf("DEBUG: About to create post in DB: %+v\n", post)
	err := h.repoManager.PostRepo.CreatePost(ctx, &post)
	if err != nil {
		fmt.Printf("ERROR: Failed to create post in DB: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post", "details": err.Error()})
		return
	}

	fmt.Printf("DEBUG: Post created successfully with ID: %s\n", post.ID.Hex())

	// Cache the created post in Redis
	cacheKey := fmt.Sprintf("post:%s", post.ID.Hex())
	err = h.redisClient.Set(ctx, cacheKey, post)
	if err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to cache post: %v\n", err)
	}

	c.JSON(http.StatusCreated, post)
}

func (h *PostHandler) GetPostByID(c *gin.Context) {
	id := c.Param("id")
	cacheKey := fmt.Sprintf("post:%s", id)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Try to get from Redis first
	var post model.Post
	err := h.redisClient.Get(ctx, cacheKey, &post)
	if err == nil {
		// Increment view count asynchronously
		go h.repoManager.PostRepo.IncrementViewCount(context.Background(), id)

		c.JSON(http.StatusOK, post)
		return
	}

	// If not in cache, get from DB
	dbPost, err := h.repoManager.PostRepo.GetPostByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Increment view count
	go h.repoManager.PostRepo.IncrementViewCount(context.Background(), id)

	// Cache the result
	err = h.redisClient.Set(ctx, cacheKey, dbPost)
	if err != nil {
		fmt.Printf("Failed to cache post: %v\n", err)
	}

	c.JSON(http.StatusOK, dbPost)
}

func (h *PostHandler) UpdatePost(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("userID")

	var post model.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get existing post to verify ownership
	existing, err := h.repoManager.PostRepo.GetPostByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	if existing.AuthorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	err = h.repoManager.PostRepo.UpdatePost(ctx, id, &post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("post:%s", id)
	h.redisClient.Delete(ctx, cacheKey)

	c.JSON(http.StatusOK, gin.H{"message": "Post updated successfully"})
}

func (h *PostHandler) DeletePost(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("userID")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get existing post to verify ownership
	existing, err := h.repoManager.PostRepo.GetPostByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	if existing.AuthorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	err = h.repoManager.PostRepo.DeletePost(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("post:%s", id)
	h.redisClient.Delete(ctx, cacheKey)

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

func (h *PostHandler) GetUserPosts(c *gin.Context) {
	userID := c.Param("userId")
	limitStr := c.DefaultQuery("limit", "20")
	skipStr := c.DefaultQuery("skip", "0")

	limit, _ := strconv.ParseInt(limitStr, 10, 64)
	skip, _ := strconv.ParseInt(skipStr, 10, 64)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	posts, err := h.repoManager.PostRepo.GetPostsByUser(ctx, userID, limit, skip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get posts"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (h *PostHandler) GetFeed(c *gin.Context) {
	userID := c.GetString("userID")
	limitStr := c.DefaultQuery("limit", "20")
	skipStr := c.DefaultQuery("skip", "0")

	limit, _ := strconv.ParseInt(limitStr, 10, 64)
	skip, _ := strconv.ParseInt(skipStr, 10, 64)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Try cache first
	cacheKey := fmt.Sprintf("user:%s:feed:%d:%d", userID, limit, skip)
	var posts []model.Post
	err := h.redisClient.Get(ctx, cacheKey, &posts)
	if err == nil {
		c.JSON(http.StatusOK, posts)
		return
	}

	// Get list of users this user follows
	following, err := h.repoManager.FollowRepo.GetFollowing(ctx, userID, 1000, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get feed"})
		return
	}

	userIDs := []string{userID} // Include own posts
	for _, follow := range following {
		userIDs = append(userIDs, follow.FollowingID)
	}

	// Get feed posts
	posts, err = h.repoManager.PostRepo.GetFeed(ctx, userIDs, limit, skip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get feed"})
		return
	}

	// Cache for 15 minutes
	h.redisClient.Set(ctx, cacheKey, posts)

	c.JSON(http.StatusOK, posts)
}

func (h *PostHandler) SearchPosts(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}

	limitStr := c.DefaultQuery("limit", "20")
	skipStr := c.DefaultQuery("skip", "0")

	limit, _ := strconv.ParseInt(limitStr, 10, 64)
	skip, _ := strconv.ParseInt(skipStr, 10, 64)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	posts, err := h.repoManager.PostRepo.SearchPosts(ctx, query, limit, skip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search posts"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (h *PostHandler) GetTrendingPosts(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	limit, _ := strconv.ParseInt(limitStr, 10, 64)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Try cache first (5 min TTL)
	cacheKey := "trending:posts"
	var posts []model.Post
	err := h.redisClient.Get(ctx, cacheKey, &posts)
	if err == nil {
		c.JSON(http.StatusOK, posts)
		return
	}

	// Get from DB
	posts, err = h.repoManager.PostRepo.GetTrendingPosts(ctx, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get trending posts"})
		return
	}

	// Cache for 5 minutes
	h.redisClient.Set(ctx, cacheKey, posts)

	c.JSON(http.StatusOK, posts)
}

func (h *PostHandler) GetPostsByTag(c *gin.Context) {
	tag := c.Param("tag")
	limitStr := c.DefaultQuery("limit", "20")
	skipStr := c.DefaultQuery("skip", "0")

	limit, _ := strconv.ParseInt(limitStr, 10, 64)
	skip, _ := strconv.ParseInt(skipStr, 10, 64)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	posts, err := h.repoManager.PostRepo.GetPostsByTag(ctx, tag, limit, skip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get posts by tag"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (h *PostHandler) GetPostsByCategory(c *gin.Context) {
	category := c.Param("category")
	limitStr := c.DefaultQuery("limit", "20")
	skipStr := c.DefaultQuery("skip", "0")

	limit, _ := strconv.ParseInt(limitStr, 10, 64)
	skip, _ := strconv.ParseInt(skipStr, 10, 64)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	posts, err := h.repoManager.PostRepo.GetPostsByCategory(ctx, category, limit, skip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get posts by category"})
		return
	}

	c.JSON(http.StatusOK, posts)
}
