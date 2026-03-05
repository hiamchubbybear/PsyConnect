package http

import (
	"consultationservice/internal/newsfeed/comment/domain"
	"consultationservice/internal/newsfeed/comment/usecase"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
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
	createCommentUC *usecase.CreateCommentUseCase
	getCommentsUC   *usecase.GetCommentsUseCase
	updateCommentUC *usecase.UpdateCommentUseCase
	deleteCommentUC *usecase.DeleteCommentUseCase
	getRepliesUC    *usecase.GetRepliesUseCase
	kafka           KafkaProducer
}

func NewHandler(
	createCommentUC *usecase.CreateCommentUseCase,
	getCommentsUC *usecase.GetCommentsUseCase,
	updateCommentUC *usecase.UpdateCommentUseCase,
	deleteCommentUC *usecase.DeleteCommentUseCase,
	getRepliesUC *usecase.GetRepliesUseCase,
	kafka KafkaProducer,
) *Handler {
	return &Handler{
		createCommentUC: createCommentUC,
		getCommentsUC:   getCommentsUC,
		updateCommentUC: updateCommentUC,
		deleteCommentUC: deleteCommentUC,
		getRepliesUC:    getRepliesUC,
		kafka:           kafka,
	}
}

func (h *Handler) CreateComment(c *gin.Context) {
	postID := c.Param("id")
	userID := c.GetString("userID")

	var req struct {
		Content         string `json:"content" binding:"required"`
		ParentCommentID string `json:"parent_comment_id,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content is required"})
		return
	}

	
	if len(req.Content) < 1 || len(req.Content) > 500 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content must be 1-500 characters"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	comment, err := h.createCommentUC.Execute(ctx, postID, userID, req.Content, req.ParentCommentID)
	if err != nil {
		if err.Error() == "max depth exceeded" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot reply more than 3 levels deep"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	
	if req.ParentCommentID == "" {
		go h.sendCommentNotification(postID, userID, comment.ID.Hex())
	}

	c.JSON(http.StatusCreated, comment)
}

func (h *Handler) GetComments(c *gin.Context) {
	postID := c.Param("id")
	limitStr := c.DefaultQuery("limit", "20")
	skipStr := c.DefaultQuery("skip", "0")
	depthStr := c.DefaultQuery("depth", "2")

	limit, _ := strconv.ParseInt(limitStr, 10, 64)
	skip, _ := strconv.ParseInt(skipStr, 10, 64)
	maxDepth, _ := strconv.Atoi(depthStr)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	comments, err := h.getCommentsUC.Execute(ctx, postID, maxDepth, limit, skip)
	if err != nil {
		log.Printf("[CommentAPI] ERROR: GetComments failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comments"})
		return
	}

	
	comments = h.populateCommentAuthors(ctx, comments)

	c.JSON(http.StatusOK, gin.H{
		"comments": comments,
		"total":    len(comments),
	})
}

func (h *Handler) UpdateComment(c *gin.Context) {
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

	err := h.updateCommentUC.Execute(ctx, commentID, userID, req.Content)
	if err != nil {
		if errors.Is(err, usecase.ErrNotAuthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment updated"})
}

func (h *Handler) DeleteComment(c *gin.Context) {
	commentID := c.Param("id")
	userID := c.GetString("userID")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := h.deleteCommentUC.Execute(ctx, commentID, userID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotAuthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted"})
}

func (h *Handler) sendCommentNotification(postID, commenterID, commentID string) {
	notificationData := map[string]interface{}{
		"postId":      postID,
		"commenterId": commenterID,
		"commentId":   commentID,
	}
	jsonData, _ := json.Marshal(notificationData)
	h.kafka.SendToTopic("notification.social.post-comment", string(jsonData))
}

func (h *Handler) populateCommentAuthors(ctx context.Context, comments []domain.Comment) []domain.Comment {
	for i := range comments {
		
		comments[i].AuthorID = comments[i].UserID

		
		if len(comments[i].Replies) > 0 {
			comments[i].Replies = h.populateCommentAuthors(ctx, comments[i].Replies)
		}
	}
	return comments
}

func (h *Handler) GetReplies(c *gin.Context) {
	parentCommentID := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	replies, err := h.getRepliesUC.Execute(ctx, parentCommentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get replies"})
		return
	}

	c.JSON(http.StatusOK, replies)
}
