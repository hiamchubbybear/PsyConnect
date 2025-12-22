package handlers

import (
	"consultationservice/bootstrap"
	"consultationservice/internal/model"
	"consultationservice/internal/redis"
	"consultationservice/internal/repository"
	"context"
	"encoding/json"
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
	userID := c.GetString("userID")

	var req struct {
		ReactionType string `json:"reaction_type" binding:"required"` // "up" or "down"
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Reaction type is required"})
		return
	}

	if req.ReactionType != model.VoteUp && req.ReactionType != model.VoteDown {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reaction type. Use 'up' or 'down'"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get existing reaction to check if we are switching
	existing, err := h.repoManager.ReactionRepo.GetUserReaction(ctx, postID, userID)

	reaction := &model.Reaction{
		PostID:       postID,
		UserID:       userID,
		ReactionType: req.ReactionType,
	}

	err = h.repoManager.ReactionRepo.AddReaction(ctx, reaction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add reaction"})
		return
	}

	// Update post engagement counts
	go func() {
		bgCtx := context.Background()
		if existing == nil {
			// New reaction
			field := "upvote_count"
			if req.ReactionType == model.VoteDown {
				field = "downvote_count"
			}
			h.repoManager.PostRepo.UpdateEngagementCount(bgCtx, postID, field, 1)
		} else if existing.ReactionType != req.ReactionType {
			// Switched reaction
			if req.ReactionType == model.VoteUp {
				h.repoManager.PostRepo.UpdateEngagementCount(bgCtx, postID, "upvote_count", 1)
				h.repoManager.PostRepo.UpdateEngagementCount(bgCtx, postID, "downvote_count", -1)
			} else {
				h.repoManager.PostRepo.UpdateEngagementCount(bgCtx, postID, "upvote_count", -1)
				h.repoManager.PostRepo.UpdateEngagementCount(bgCtx, postID, "downvote_count", 1)
			}
		}

		// Emit notification only for Upvotes
		if req.ReactionType == model.VoteUp && (existing == nil || existing.ReactionType != model.VoteUp) {
			post, err := h.repoManager.PostRepo.GetPostByID(bgCtx, postID)
			if err == nil && post != nil && post.AuthorID != userID {
				voterProfile, _ := h.repoManager.GrpcProfile.GetProfile(userID)
				voterName := "Someone"
				if voterProfile != nil {
					voterName = fmt.Sprintf("%s %s", voterProfile.FirstName, voterProfile.LastName)
				}

				notificationData := map[string]interface{}{
					"userId":    post.AuthorID,
					"voterName": voterName,
					"postId":    postID,
				}
				jsonData, _ := json.Marshal(notificationData)
				h.repoManager.Kafka.SendToTopic("notification.social.post-upvote", string(jsonData))
			}
		}
	}()

	c.JSON(http.StatusOK, reaction)
}

func (h *ReactionHandler) RemoveReaction(c *gin.Context) {
	postID := c.Param("id")
	userID := c.GetString("userID")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get existing reaction to know which count to decrement
	existing, err := h.repoManager.ReactionRepo.GetUserReaction(ctx, postID, userID)
	if err != nil || existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reaction not found"})
		return
	}

	err = h.repoManager.ReactionRepo.RemoveReaction(ctx, postID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove reaction"})
		return
	}

	// Update post engagement count
	go func() {
		field := "upvote_count"
		if existing.ReactionType == model.VoteDown {
			field = "downvote_count"
		}
		h.repoManager.PostRepo.UpdateEngagementCount(context.Background(), postID, field, -1)
	}()

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
