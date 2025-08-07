package handlers

import (
	"chatservice/bootstrap"
	"chatservice/internal/model"
	"chatservice/internal/repository"
	"chatservice/pkg/apiresponse"
	"net/http"

	"github.com/gin-gonic/gin"
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

	chats, err := h.RepoManager.MessageRepo.FindChatsByConversation(conversationID)
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
