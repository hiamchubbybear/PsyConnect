package http

import (
	"consultationservice/internal/newsfeed/reaction/domain"
	"consultationservice/internal/newsfeed/reaction/usecase"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)


type KafkaProducer interface {
	SendToTopic(topic string, message string) error
}


type ProfileService interface {
	GetProfile(userID string) (interface{}, error)
}

type Handler struct {
	toggleReactionUC *usecase.ToggleReactionUseCase
	getReactionsUC   *usecase.GetReactionsUseCase
	kafka            KafkaProducer
}

func NewHandler(
	toggleReactionUC *usecase.ToggleReactionUseCase,
	getReactionsUC *usecase.GetReactionsUseCase,
	kafka KafkaProducer,
) *Handler {
	return &Handler{
		toggleReactionUC: toggleReactionUC,
		getReactionsUC:   getReactionsUC,
		kafka:            kafka,
	}
}

func (h *Handler) AddReaction(c *gin.Context) {
	postID := c.Param("id")
	userID := c.GetString("userID")

	var req struct {
		ReactionType string `json:"reaction_type" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Reaction type is required"})
		return
	}

	if req.ReactionType != domain.VoteUp && req.ReactionType != domain.VoteDown {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reaction type. Use 'up' or 'down'"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	existing, newReaction, err := h.toggleReactionUC.Execute(ctx, postID, userID, req.ReactionType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle reaction"})
		return
	}

	
	if newReaction != nil && newReaction.IsUpvote() && (existing == nil || !existing.IsUpvote()) {
		go h.sendUpvoteNotification(postID, userID)
	}

	if newReaction == nil {
		
		c.JSON(http.StatusOK, gin.H{"message": "Reaction removed"})
		return
	}

	c.JSON(http.StatusOK, newReaction)
}

func (h *Handler) RemoveReaction(c *gin.Context) {
	postID := c.Param("id")
	userID := c.GetString("userID")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	
	existing, _, err := h.toggleReactionUC.Execute(ctx, postID, userID, "")
	if err != nil || existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reaction not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reaction removed"})
}

func (h *Handler) GetPostReactions(c *gin.Context) {
	postID := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	reactions, err := h.getReactionsUC.Execute(ctx, postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get reactions"})
		return
	}

	c.JSON(http.StatusOK, reactions)
}

func (h *Handler) sendUpvoteNotification(postID, voterID string) {
	
	
	notificationData := map[string]interface{}{
		"postId":  postID,
		"voterId": voterID,
	}
	jsonData, _ := json.Marshal(notificationData)
	h.kafka.SendToTopic("notification.social.post-upvote", string(jsonData))
}
