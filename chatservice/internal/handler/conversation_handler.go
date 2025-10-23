package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"chatservice/bootstrap"
	"chatservice/internal/model"
	"chatservice/internal/repository"
	encoder "chatservice/internal/utils/conversation"
)

type ConversationHandler struct {
	RepoManager *repository.RepositoryManager
}

func NewConversationHandler(env *bootstrap.Env, repoManager *repository.RepositoryManager) *ConversationHandler {
	return &ConversationHandler{
		RepoManager: repoManager,
	}
}

func (h *ConversationHandler) GetConversationByUsers(c *gin.Context) {
	user1 := c.Query("user1")
	user2 := c.Query("user2")

	if user1 == "" || user2 == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "missing user ids"})
		return
	}

	conv, err := h.RepoManager.ConversationRepository.GetConversationByUsers(user1, user2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get conversation"})
		return
	}
	conversationID, err := encoder.New().EncodeConversationId(user1, user2)
	if err != nil {
		c.JSON(http.StatusUpgradeRequired, gin.H{"message": "failed to encode conversation"})
	}
	if conv == nil {
		newConv := &model.Conversation{
			Id:           conversationID,
			Participants: []string{user1, user2},
			CreatedAt:    model.NowUTC(),
		}

		err := h.RepoManager.ConversationRepository.CreateConversation(newConv)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create conversation"})
			return
		}

		log.Printf("[ConversationHandler] Created new conversation between %s and %s", user1, user2)
		c.JSON(http.StatusCreated, gin.H{
			"message": "conversation created",
			"data":    newConv,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    conv,
	})
}
