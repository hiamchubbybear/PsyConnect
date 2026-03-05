package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"notificationservice/pkg/email"
	"notificationservice/pkg/handlers"
)

type API struct {
	emailSvc *email.EmailService
	notifSvc *handlers.NotificationService
}

func NewAPI(emailSvc *email.EmailService, notifSvc *handlers.NotificationService) *API {
	return &API{
		emailSvc: emailSvc,
		notifSvc: notifSvc,
	}
}



type ActivateEmailRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Code     string `json:"code" binding:"required"`
	Fullname string `json:"fullname" binding:"required"`
}

func (a *API) SendActivateEmail(c *gin.Context) {
	var req ActivateEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := a.emailSvc.SendActivationEmail(req.Email, req.Username, req.Code, req.Fullname); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Activation email sent successfully"})
}

type AccountUpdateEmailRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

func (a *API) SendAccountUpdateEmail(c *gin.Context) {
	var req AccountUpdateEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := a.emailSvc.SendAccountChangeEmail(req.Email, req.Username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account update email sent successfully"})
}



type SaveTokenRequest struct {
	UserID string `json:"userId" binding:"required"`
	Token  string `json:"token" binding:"required"`
}

func (a *API) SaveFCMToken(c *gin.Context) {
	var req SaveTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := a.notifSvc.SaveFCMToken(req.UserID, req.Token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "FCM token saved successfully"})
}

func (a *API) GetNotifications(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	skip, _ := strconv.Atoi(c.DefaultQuery("skip", "0"))

	notifications, err := a.notifSvc.GetNotifications(userID, limit, skip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get notifications"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"notifications": notifications,
		"count":         len(notifications),
	})
}

type MarkAsReadRequest struct {
	UserID string `json:"userId" binding:"required"`
}

func (a *API) MarkAsRead(c *gin.Context) {
	notifIDStr := c.Param("id")
	notifID, err := strconv.ParseUint(notifIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	var req MarkAsReadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := a.notifSvc.MarkAsRead(uint(notifID), req.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}
