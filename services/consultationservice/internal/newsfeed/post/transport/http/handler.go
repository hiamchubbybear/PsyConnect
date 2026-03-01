package http

import (
	"consultationservice/internal/newsfeed/post/domain"
	"consultationservice/internal/newsfeed/post/usecase"
	"consultationservice/internal/redis"
	"consultationservice/internal/repository"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	createPostUC     *usecase.CreatePostUseCase
	getPostByIDUC    *usecase.GetPostByIDUseCase
	updatePostUC     *usecase.UpdatePostUseCase
	deletePostUC     *usecase.DeletePostUseCase
	getFeedUC        *usecase.GetFeedUseCase
	getTrendingUC    *usecase.GetTrendingPostsUseCase
	searchPostsUC    *usecase.SearchPostsUseCase
	getUserPostsUC   *usecase.GetUserPostsUseCase
	getByTagUC       *usecase.GetPostsByTagUseCase
	getByCategoryUC  *usecase.GetPostsByCategoryUseCase
	incrementViewUC  *usecase.IncrementViewCountUseCase
	getPopularTagsUC *usecase.GetPopularTagsUseCase
	redisClient      redis.RedisStore
	reactionRepo     repository.ReactionRepository
	bookmarkRepo     repository.BookmarkRepository
}

func NewHandler(
	createPostUC *usecase.CreatePostUseCase,
	getPostByIDUC *usecase.GetPostByIDUseCase,
	updatePostUC *usecase.UpdatePostUseCase,
	deletePostUC *usecase.DeletePostUseCase,
	getFeedUC *usecase.GetFeedUseCase,
	getTrendingUC *usecase.GetTrendingPostsUseCase,
	searchPostsUC *usecase.SearchPostsUseCase,
	getUserPostsUC *usecase.GetUserPostsUseCase,
	getByTagUC *usecase.GetPostsByTagUseCase,
	getByCategoryUC *usecase.GetPostsByCategoryUseCase,
	incrementViewUC *usecase.IncrementViewCountUseCase,
	getPopularTagsUC *usecase.GetPopularTagsUseCase,
	redisClient redis.RedisStore,
	reactionRepo repository.ReactionRepository,
	bookmarkRepo repository.BookmarkRepository,
) *Handler {
	return &Handler{
		createPostUC:     createPostUC,
		getPostByIDUC:    getPostByIDUC,
		updatePostUC:     updatePostUC,
		deletePostUC:     deletePostUC,
		getFeedUC:        getFeedUC,
		getTrendingUC:    getTrendingUC,
		searchPostsUC:    searchPostsUC,
		getUserPostsUC:   getUserPostsUC,
		getByTagUC:       getByTagUC,
		getByCategoryUC:  getByCategoryUC,
		incrementViewUC:  incrementViewUC,
		getPopularTagsUC: getPopularTagsUC,
		redisClient:      redisClient,
		reactionRepo:     reactionRepo,
		bookmarkRepo:     bookmarkRepo,
	}
}

func (h *Handler) CreatePost(c *gin.Context) {
	var post domain.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authorID := c.GetString("userID")
	post.AuthorID = authorID

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := h.createPostUC.Execute(ctx, &post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	// Cache the created post
	cacheKey := fmt.Sprintf("post:%s", post.ID.Hex())
	h.redisClient.Set(ctx, cacheKey, post)

	c.JSON(http.StatusCreated, post)
}

func (h *Handler) GetPostByID(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("userID")
	cacheKey := fmt.Sprintf("post:%s", id)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Try Redis cache first
	var post domain.Post
	err := h.redisClient.Get(ctx, cacheKey, &post)
	if err == nil {
		h.populateUserState(ctx, &post, userID)
		go h.incrementViewUC.Execute(context.Background(), id)
		c.JSON(http.StatusOK, post)
		return
	}

	// Get from DB
	dbPost, err := h.getPostByIDUC.Execute(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	h.populateUserState(ctx, dbPost, userID)
	go h.incrementViewUC.Execute(context.Background(), id)
	h.redisClient.Set(ctx, cacheKey, dbPost)

	c.JSON(http.StatusOK, dbPost)
}

func (h *Handler) UpdatePost(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("userID")

	var updates domain.Post
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := h.updatePostUC.Execute(ctx, id, userID, &updates)
	if err != nil {
		if err == usecase.ErrNotAuthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("post:%s", id)
	h.redisClient.Delete(ctx, cacheKey)

	c.JSON(http.StatusOK, gin.H{"message": "Post updated successfully"})
}

func (h *Handler) DeletePost(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("userID")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := h.deletePostUC.Execute(ctx, id, userID)
	if err != nil {
		if err == usecase.ErrNotAuthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("post:%s", id)
	h.redisClient.Delete(ctx, cacheKey)

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

func (h *Handler) GetFeed(c *gin.Context) {
	userID := c.GetString("userID")
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64)
	skip, _ := strconv.ParseInt(c.DefaultQuery("skip", "0"), 10, 64)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var excludeIDs []string
	if userID != "" {
		viewedKey := fmt.Sprintf("user:%s:viewed", userID)
		ids, err := h.redisClient.SMembers(ctx, viewedKey)
		if err == nil {
			excludeIDs = ids
		}
	}

	posts, err := h.getFeedUC.Execute(ctx, []string{userID}, excludeIDs, limit, skip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get feed"})
		return
	}

	posts = h.populateUserStates(ctx, posts, userID)
	c.JSON(http.StatusOK, posts)
}

func (h *Handler) GetTrendingPosts(c *gin.Context) {
	userID := c.GetString("userID")
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	posts, err := h.getTrendingUC.Execute(ctx, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get trending posts"})
		return
	}

	posts = h.populateUserStates(ctx, posts, userID)
	c.JSON(http.StatusOK, posts)
}

func (h *Handler) SearchPosts(c *gin.Context) {
	userID := c.GetString("userID")
	query := c.Query("q")
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64)
	skip, _ := strconv.ParseInt(c.DefaultQuery("skip", "0"), 10, 64)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	posts, err := h.searchPostsUC.Execute(ctx, query, limit, skip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search posts"})
		return
	}

	posts = h.populateUserStates(ctx, posts, userID)
	c.JSON(http.StatusOK, posts)
}

func (h *Handler) GetUserPosts(c *gin.Context) {
	currentUserID := c.GetString("userID")
	targetUserID := c.Param("userId")
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64)
	skip, _ := strconv.ParseInt(c.DefaultQuery("skip", "0"), 10, 64)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	posts, err := h.getUserPostsUC.Execute(ctx, targetUserID, limit, skip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user posts"})
		return
	}

	posts = h.populateUserStates(ctx, posts, currentUserID)
	c.JSON(http.StatusOK, posts)
}

func (h *Handler) GetPostsByTag(c *gin.Context) {
	userID := c.GetString("userID")
	tag := c.Param("tag")
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64)
	skip, _ := strconv.ParseInt(c.DefaultQuery("skip", "0"), 10, 64)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	posts, err := h.getByTagUC.Execute(ctx, tag, limit, skip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get posts by tag"})
		return
	}

	posts = h.populateUserStates(ctx, posts, userID)
	c.JSON(http.StatusOK, posts)
}

func (h *Handler) GetPostsByCategory(c *gin.Context) {
	userID := c.GetString("userID")
	category := c.Param("category")
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64)
	skip, _ := strconv.ParseInt(c.DefaultQuery("skip", "0"), 10, 64)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	posts, err := h.getByCategoryUC.Execute(ctx, category, limit, skip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get posts by category"})
		return
	}

	posts = h.populateUserStates(ctx, posts, userID)
	c.JSON(http.StatusOK, posts)
}

func (h *Handler) GetPopularTags(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(limitStr)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tags, err := h.getPopularTagsUC.Execute(ctx, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get popular tags"})
		return
	}

	c.JSON(http.StatusOK, tags)
}

func (h *Handler) IncrementViewCount(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("userID")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Track in DB
	err := h.incrementViewUC.Execute(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to increment view count"})
		return
	}

	// Track in Redis for "Hide Viewed Posts"
	if userID != "" {
		viewedKey := fmt.Sprintf("user:%s:viewed", userID)
		h.redisClient.SAdd(ctx, viewedKey, id)
	}

	c.JSON(http.StatusOK, gin.H{"message": "View recorded"})
}

// Helper methods for user state population
func (h *Handler) populateUserState(ctx context.Context, post *domain.Post, userID string) {
	if userID == "" {
		return
	}

	// Get user's reaction
	reaction, err := h.reactionRepo.GetUserReaction(ctx, post.ID.Hex(), userID)
	if err == nil && reaction != nil {
		post.UserVote = &reaction.ReactionType
	}

	// Get user's bookmark status
	isBookmarked, err := h.bookmarkRepo.IsBookmarked(ctx, userID, post.ID.Hex())
	if err == nil {
		post.UserBookmark = isBookmarked
	}
}

func (h *Handler) populateUserStates(ctx context.Context, posts []domain.Post, userID string) []domain.Post {
	if userID == "" {
		return posts
	}

	for i := range posts {
		h.populateUserState(ctx, &posts[i], userID)
	}
	return posts
}
