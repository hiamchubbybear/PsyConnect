package handler

import (
	"chatservice/internal/chat/model/model"
	"chatservice/internal/chat/repository/repository"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"chatservice/bootstrap"
	encoder "chatservice/internal/utils/conversation"
	"sort"
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
			ID:           conversationID,
			Participants: []string{user1, user2},
			CreatedAt:    model.NowUTC(),
		}

		err := h.RepoManager.ConversationRepository.CreateConversation(newConv)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create conversation"})
			return
		}

		// Inject system welcome message into the new conversation
		systemMsg := &model.Chat{
			ID:             model.NewUUID(),
			SenderID:       "SYSTEM",
			ConversationID: conversationID,
			Text:           "🤝 Cuộc trò chuyện đã bắt đầu. Hãy tôn trọng và lịch sự với nhau. PsyConnect khuyến khích bạn giữ thái độ tích cực trong mọi cuộc trò chuyện.",
			IsSystem:       true,
			CreatedAt:      time.Now().UTC(),
			UpdatedAt:      time.Now().UTC(),
		}
		if _, sysErr := h.RepoManager.MessageRepo.CreateChat(systemMsg); sysErr != nil {
			log.Printf("[ConversationHandler] Warning: failed to insert system message: %v", sysErr)
		} else {
			log.Printf("[ConversationHandler] System welcome message sent to conversation %s", conversationID)
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

type RecentConversationResponse struct {
	*model.Conversation
	LastMessage *model.Chat `json:"lastMessage"`
}

func (h *ConversationHandler) GetRecentConversations(c *gin.Context) {
	userID := c.Query("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "missing userId"})
		return
	}

	conversations, err := h.RepoManager.ConversationRepository.GetConversationsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to limit conversations"})
		return
	}

	var dtos []RecentConversationResponse
	for _, conv := range conversations {
		dto := RecentConversationResponse{
			Conversation: conv,
		}

		lastMsg, _ := h.RepoManager.MessageRepo.FindLastChatByConversation(conv.ID)
		dto.LastMessage = lastMsg

		dtos = append(dtos, dto)
	}

	// Sort by last message timestamp descending
	sort.Slice(dtos, func(i, j int) bool {
		var timeI, timeJ time.Time
		if dtos[i].LastMessage != nil {
			timeI = dtos[i].LastMessage.CreatedAt
		} else {
			timeI = dtos[i].Conversation.CreatedAt
		}

		if dtos[j].LastMessage != nil {
			timeJ = dtos[j].LastMessage.CreatedAt
		} else {
			timeJ = dtos[j].Conversation.CreatedAt
		}

		return timeI.After(timeJ)
	})

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    dtos,
	})
}
