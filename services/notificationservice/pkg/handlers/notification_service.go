package handlers

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/psyconnect/notificationservice/pkg/database"
	"github.com/psyconnect/notificationservice/pkg/firebase"
)

type NotificationService struct {
	db             *database.Database
	fcm            *firebase.FCMService
	spamPrevention *SpamPrevention
}

// SpamPrevention prevents sending duplicate notifications too quickly
type SpamPrevention struct {
	mu           sync.RWMutex
	recentNotifs map[string]time.Time
	spamWindowMs int64
}

func NewSpamPrevention(windowMs int64) *SpamPrevention {
	return &SpamPrevention{
		recentNotifs: make(map[string]time.Time),
		spamWindowMs: windowMs,
	}
}

func (sp *SpamPrevention) CanSend(userID, notifType string) bool {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	key := userID + ":" + notifType
	now := time.Now()

	if lastSent, exists := sp.recentNotifs[key]; exists {
		if now.Sub(lastSent).Milliseconds() < sp.spamWindowMs {
			return false
		}
	}

	sp.recentNotifs[key] = now
	return true
}

func NewNotificationService(db *database.Database, fcm *firebase.FCMService) *NotificationService {
	return &NotificationService{
		db:             db,
		fcm:            fcm,
		spamPrevention: NewSpamPrevention(60000), // 60 seconds
	}
}

// SendToUser sends a notification to a specific user
func (ns *NotificationService) SendToUser(userID, title, body, notifType string, metadata map[string]interface{}) error {
	// Save to database
	metadataJSON, _ := json.Marshal(metadata)
	notification := &database.Notification{
		UserID:   userID,
		Title:    title,
		Body:     body,
		Type:     notifType,
		Metadata: string(metadataJSON),
		IsRead:   false,
	}

	if err := ns.db.SaveNotification(notification); err != nil {
		log.Printf("❌ Failed to save notification: %v", err)
		return err
	}

	log.Printf("✅ Notification saved to database for user: %s", userID)

	// Try to send push notification
	if ns.fcm != nil {
		token, err := ns.db.GetFCMToken(userID)
		if err != nil {
			log.Printf("⚠️  No FCM token for user %s, notification only persisted", userID)
			return nil
		}

		// Convert metadata to string map for FCM
		dataMap := make(map[string]string)
		for k, v := range metadata {
			if str, ok := v.(string); ok {
				dataMap[k] = str
			}
		}
		dataMap["type"] = notifType
		dataMap["notificationId"] = string(rune(notification.ID))

		if err := ns.fcm.SendPushNotification(token, title, body, dataMap); err != nil {
			log.Printf("❌ Failed to send push notification: %v", err)
			return err
		}

		log.Printf("✅ Push notification sent to user: %s", userID)
	}

	return nil
}

// SendToUserWithSpamCheck sends notification with spam prevention
func (ns *NotificationService) SendToUserWithSpamCheck(userID, title, body, notifType string, metadata map[string]interface{}) error {
	if !ns.spamPrevention.CanSend(userID, notifType) {
		log.Printf("⚠️  Spam prevention: skipping notification for user %s, type %s", userID, notifType)
		return nil
	}

	return ns.SendToUser(userID, title, body, notifType, metadata)
}

// GetNotifications retrieves notifications for a user
func (ns *NotificationService) GetNotifications(userID string, limit, offset int) ([]database.Notification, error) {
	return ns.db.GetNotifications(userID, limit, offset)
}

// MarkAsRead marks a notification as read
func (ns *NotificationService) MarkAsRead(notificationID uint, userID string) error {
	return ns.db.MarkAsRead(notificationID, userID)
}

// SaveFCMToken saves FCM token for a user
func (ns *NotificationService) SaveFCMToken(userID, token string) error {
	return ns.db.SaveFCMToken(userID, token)
}
