package handler

import (
	"chatservice/internal/chat/model/model"
	"chatservice/internal/chat/repository/repository"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"chatservice/bootstrap"

	encoder "chatservice/internal/utils/conversation"
	"chatservice/pkg/apiresponse"
)

type ChatHandler struct {
	RepoManager *repository.RepositoryManager
}

func NewChatHandler(env *bootstrap.Env, repoManager *repository.RepositoryManager) *ChatHandler {
	return &ChatHandler{
		RepoManager: repoManager,
	}
}

func (h *ChatHandler) CreateChat(c *gin.Context) {
	var chat model.Chat

	if err := c.ShouldBindJSON(&chat); err != nil {
		apiresponse.ErrorHandler(c, http.StatusBadRequest, "Invalid input")
		return
	}

	if chat.ConversationID == "" || chat.SenderID == "" {
		apiresponse.ErrorHandler(c, http.StatusBadRequest, "Missing ConversationID or UserID")
		return
	}

	chat.ID = model.NewUUID()
	chat.CreatedAt = model.NowUTC()

	id, err := h.RepoManager.MessageRepo.CreateChat(&chat)
	if err != nil {
		apiresponse.ErrorHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	apiresponse.NewApiResponse(c, gin.H{"id": id, "message": chat})
}

func (h *ChatHandler) GetChatByID(c *gin.Context) {
	chatID := c.Param("id")
	if chatID == "" {
		apiresponse.ErrorHandler(c, http.StatusBadRequest, "Missing chat ID")
		return
	}

	chat, err := h.RepoManager.MessageRepo.FindChatByID(chatID)
	if err != nil {
		apiresponse.ErrorHandler(c, http.StatusInternalServerError, err.Error())
		return
	}
	if chat == nil {
		apiresponse.ErrorHandler(c, http.StatusNotFound, "Chat not found")
		return
	}

	apiresponse.NewApiResponse(c, chat)
}

func (h *ChatHandler) GetChatsByConversation(c *gin.Context) {
	conversationID := c.Param("conversationId")
	if conversationID == "" {
		apiresponse.ErrorHandler(c, http.StatusBadRequest, "Missing conversation ID")
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	beforeStr := c.Query("before")
	var before time.Time
	if beforeStr != "" {
		before, err = time.Parse(time.RFC3339, beforeStr)
		if err != nil {
			apiresponse.ErrorHandler(c, http.StatusBadRequest, "Invalid 'before' time format, must be RFC3339")
			return
		}
	}
	chats, err := h.RepoManager.MessageRepo.FindChatsByConversation(conversationID, limit, before)
	if err != nil {
		apiresponse.ErrorHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	apiresponse.NewApiResponse(c, chats)
}

func (h *ChatHandler) UpdateChat(c *gin.Context) {
	chatID := c.Param("id")
	if chatID == "" {
		apiresponse.ErrorHandler(c, http.StatusBadRequest, "Missing chat ID")
		return
	}

	var req struct {
		Text string `json:"text"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.Text == "" {
		apiresponse.ErrorHandler(c, http.StatusBadRequest, "Invalid input")
		return
	}

	updated, err := h.RepoManager.MessageRepo.UpdateChatByID(chatID, req.Text)
	if err != nil {
		apiresponse.ErrorHandler(c, http.StatusInternalServerError, err.Error())
		return
	}
	if updated == nil {
		apiresponse.ErrorHandler(c, http.StatusNotFound, "Chat not found")
		return
	}

	apiresponse.NewApiResponse(c, updated)
}

func (h *ChatHandler) DeleteChat(c *gin.Context) {
	chatID := c.Param("id")
	if chatID == "" {
		apiresponse.ErrorHandler(c, http.StatusBadRequest, "Missing chat ID")
		return
	}

	count, err := h.RepoManager.MessageRepo.DeleteChatByID(chatID)
	if err != nil {
		apiresponse.ErrorHandler(c, http.StatusInternalServerError, err.Error())
		return
	}
	if count == 0 {
		apiresponse.ErrorHandler(c, http.StatusNotFound, "Chat not found")
		return
	}

	apiresponse.NewApiResponse(c, gin.H{"deleted": true})
}
func (h *ChatHandler) CreateConversation(c *gin.Context) {

	var req struct {
		UserIds []string `json:"userIds"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.UserIds) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "need at least two user IDs"})
		return
	}
	fmt.Println("Users " + req.UserIds[0])
	conversationUUID, err := encoder.New().EncodeConversationId(req.UserIds[0], req.UserIds[1])
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	existing, _ := h.RepoManager.ConversationRepository.Exists(conversationUUID)
	if existing {
		c.JSON(http.StatusOK, gin.H{"data": existing})
		return
	}
	conversation := model.Conversation{
		ID:           conversationUUID,
		Participants: req.UserIds,
		CreatedAt:    time.Now(),
	}

	if err := h.RepoManager.ConversationRepository.CreateConversation(&conversation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": conversation})
}
func (s *ChatHandler) GetOrCreateConversation(user1ID, user2ID string) (*model.Conversation, error) {
	existing, err := s.RepoManager.ConversationRepository.GetConversationByUsers(user1ID, user2ID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}

	newConv := &model.Conversation{
		Participants: []string{user1ID, user2ID},
	}
	err = s.RepoManager.ConversationRepository.CreateConversation(newConv)
	if err != nil {
		return nil, err
	}
	return newConv, nil
}
