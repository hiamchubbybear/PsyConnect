package handlers

import (
	"consultationservice/bootstrap"
	"consultationservice/internal/model"
	"consultationservice/internal/redis"
	"consultationservice/internal/repository"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ReactionHandler struct {
	env         *bootstrap.Env
	repoManager *repository.RepositoryManager
	redisClient redis.RedisStore
}

func NewReactionHandler(env *bootstrap.Env, repoManager *repository.RepositoryManager, redisClient redis.RedisStore) *ReactionHandler {
	return &ReactionHandler{
		env:         env,
		repoManager: repoManager,
		redisClient: redisClient,
	}
}

func (h *ReactionHandler) AddReaction(c *gin.Context) {
	postID := c.Param("id")
	userID := c.GetString("userID") // From auth middleware

	var req struct {
		ReactionType string `json:"reaction_type" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate reaction type
	validReactions := map[string]bool{
		model.ReactionLike:  true,
		model.ReactionLove:  true,
		model.ReactionLaugh: true,
		model.ReactionThink: true,
		model.ReactionSad:   true,
		model.ReactionAngry: true,
	}

	if !validReactions[req.ReactionType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid reaction type"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	reaction := &model.Reaction{
		PostID:       postID,
		UserID:       userID,
		ReactionType: req.ReactionType,
	}

	err := h.repoManager.ReactionRepo.AddReaction(ctx, reaction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add reaction"})
		return
	}

	// Update post like count
	err = h.repoManager.PostRepo.UpdateEngagementCount(ctx, postID, "like_count", 1)
	if err != nil {
		fmt.Printf("Failed to update like count: %v\n", err)
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("post:%s", postID)
	h.redisClient.Delete(ctx, cacheKey)

	c.JSON(http.StatusCreated, reaction)
}

func (h *ReactionHandler) RemoveReaction(c *gin.Context) {
	postID := c.Param("id")
	userID := c.GetString("userID")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := h.repoManager.ReactionRepo.RemoveReaction(ctx, postID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reaction not found"})
		return
	}

	// Update post like count
	err = h.repoManager.PostRepo.UpdateEngagementCount(ctx, postID, "like_count", -1)
	if err != nil {
		fmt.Printf("Failed to update like count: %v\n", err)
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("post:%s", postID)
	h.redisClient.Delete(ctx, cacheKey)

	c.JSON(http.StatusOK, gin.H{"message": "Reaction removed"})
}

func (h *ReactionHandler) GetPostReactions(c *gin.Context) {
	postID := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Try cache first
	cacheKey := fmt.Sprintf("post:%s:reactions", postID)
	var reactions []model.Reaction
	err := h.redisClient.Get(ctx, cacheKey, &reactions)
	if err == nil {
		c.JSON(http.StatusOK, reactions)
		return
	}

	// Get from DB
	reactions, err = h.repoManager.ReactionRepo.GetReactionsByPost(ctx, postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get reactions"})
		return
	}

	// Cache result
	h.redisClient.Set(ctx, cacheKey, reactions)

	c.JSON(http.StatusOK, reactions)
}
