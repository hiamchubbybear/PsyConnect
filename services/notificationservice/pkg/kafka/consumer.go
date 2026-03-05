package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"notificationservice/pkg/config"
	"notificationservice/pkg/email"
	"notificationservice/pkg/handlers"
	"notificationservice/pkg/models"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader   *kafka.Reader
	emailSvc *email.EmailService
	notifSvc *handlers.NotificationService
	config   *config.Config
}

func NewConsumer(cfg *config.Config, emailSvc *email.EmailService, notifSvc *handlers.NotificationService) *Consumer {
	
	topics := []string{
		
		"notification.user-create",
		"notification.user-activate",
		"notification.account-change",
		"notification.user-reset",
		
		"notification.therapist-approve",
		"notification.therapist-reject",
		
		"notification.push.new-message",
		"notification.push.consultation-created",
		"notification.push.consultation-updated",
		"notification.push.consultation-reminder",
		"notification.push.consultation-completed",
		"notification.push.profile-update",
		"notification.push.profile-approved",
		"notification.push.profile-rejected",
		"notification.push.new-review",
		"notification.push.new-client",
		"notification.push.system",
		
		"notification.social.post-upvote",
		"notification.social.post-bookmark",
		"notification.social.post-share",
		"notification.social.post-comment",
		"notification.social.user-follow",
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        cfg.KafkaBrokers,
		GroupID:        cfg.KafkaGroupID,
		GroupTopics:    topics,
		MinBytes:       10e3,
		MaxBytes:       10e6,
		CommitInterval: time.Second,
		StartOffset:    kafka.LastOffset,
	})

	return &Consumer{
		reader:   reader,
		emailSvc: emailSvc,
		notifSvc: notifSvc,
		config:   cfg,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	log.Println("🚀 Starting Kafka consumer...")
	log.Printf("📡 Subscribed to 22 topics")
	log.Printf("👥 Consumer group: %s", c.config.KafkaGroupID)

	for {
		select {
		case <-ctx.Done():
			log.Println("🛑 Shutting down Kafka consumer...")
			return c.reader.Close()
		default:
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				if err == context.Canceled {
					return nil
				}
				log.Printf("❌ Error reading message: %v", err)
				continue
			}

			log.Printf("📨 Received message from topic: %s", msg.Topic)

			if err := c.processMessage(msg); err != nil {
				log.Printf("❌ Error processing message: %v", err)
			}
		}
	}
}

func (c *Consumer) processMessage(msg kafka.Message) error {
	topic := msg.Topic

	switch topic {
	
	case "notification.user-create":
		return c.handleUserCreate(msg.Value)
	case "notification.user-activate":
		return c.handleUserActivate(msg.Value)
	case "notification.user-reset":
		return c.handlePasswordReset(msg.Value)
	case "notification.account-change":
		return c.handleAccountChange(msg.Value)

	
	case "notification.therapist-approve":
		return c.handleTherapistApprove(msg.Value)
	case "notification.therapist-reject":
		return c.handleTherapistReject(msg.Value)

	
	case "notification.push.new-message":
		return c.handleNewMessage(msg.Value)
	case "notification.push.consultation-created":
		return c.handleConsultationCreated(msg.Value)
	case "notification.push.consultation-updated":
		return c.handleConsultationUpdated(msg.Value)
	case "notification.push.consultation-reminder":
		return c.handleConsultationReminder(msg.Value)
	case "notification.push.consultation-completed":
		return c.handleConsultationCompleted(msg.Value)
	case "notification.push.profile-update":
		return c.handleProfileUpdate(msg.Value)
	case "notification.push.profile-approved":
		return c.handleProfileApproved(msg.Value)
	case "notification.push.profile-rejected":
		return c.handleProfileRejected(msg.Value)
	case "notification.push.new-review":
		return c.handleNewReview(msg.Value)
	case "notification.push.new-client":
		return c.handleNewClient(msg.Value)
	case "notification.push.system":
		return c.handleSystem(msg.Value)

	case "consultation.incoming_call":
		return c.handleIncomingCall(msg.Value)

	
	case "notification.social.post-upvote":
		return c.handlePostUpvote(msg.Value)
	case "notification.social.post-bookmark":
		return c.handlePostBookmark(msg.Value)
	case "notification.social.post-share":
		return c.handlePostShare(msg.Value)
	case "notification.social.post-comment":
		return c.handlePostComment(msg.Value)
	case "notification.social.user-follow":
		return c.handleUserFollow(msg.Value)

	default:
		log.Printf("⚠️  Unknown topic: %s", topic)
		return nil
	}
}


func (c *Consumer) handleUserCreate(data []byte) error {
	var notif models.UserCreateNotification
	if err := json.Unmarshal(data, &notif); err != nil {
		return err
	}
	log.Printf("📧 Sending activation email to: %s", notif.Email)
	return c.emailSvc.SendActivationEmail(notif.Email, notif.Username, notif.Code, notif.Fullname)
}

func (c *Consumer) handleUserActivate(data []byte) error {
	var notif models.UserActivateNotification
	if err := json.Unmarshal(data, &notif); err != nil {
		return err
	}
	log.Printf("📧 Sending re-activation email to: %s", notif.Email)
	return c.emailSvc.SendActivationEmail(notif.Email, notif.Username, notif.Code, notif.Fullname)
}

func (c *Consumer) handlePasswordReset(data []byte) error {
	var notif models.PasswordResetNotification
	if err := json.Unmarshal(data, &notif); err != nil {
		return err
	}
	log.Printf("🔐 Sending password reset email to: %s", notif.Email)
	return c.emailSvc.SendPasswordResetEmail(notif.Email, notif.Username, notif.Code)
}

func (c *Consumer) handleAccountChange(data []byte) error {
	var notif models.AccountChangeNotification
	if err := json.Unmarshal(data, &notif); err != nil {
		return err
	}
	log.Printf("📝 Sending account change email to: %s", notif.Email)
	return c.emailSvc.SendAccountChangeEmail(notif.Email, notif.Username)
}


func (c *Consumer) handleTherapistApprove(data []byte) error {
	var notif models.TherapistApproveNotification
	if err := json.Unmarshal(data, &notif); err != nil {
		return err
	}
	log.Printf("✅ Therapist approved: %s", notif.Email)
	return c.notifSvc.SendToUser(
		notif.UserID,
		"Your profile was approved!",
		"Congratulations! Your therapist profile has been approved.",
		"therapist-approve",
		map[string]interface{}{},
	)
}

func (c *Consumer) handleTherapistReject(data []byte) error {
	var notif models.TherapistRejectNotification
	if err := json.Unmarshal(data, &notif); err != nil {
		return err
	}
	log.Printf("❌ Therapist rejected: %s", notif.Email)
	return c.notifSvc.SendToUser(
		notif.UserID,
		"Profile rejected",
		"Unfortunately, your profile was rejected. Please update and resubmit.",
		"therapist-reject",
		map[string]interface{}{},
	)
}


func (c *Consumer) handleNewMessage(data []byte) error {
	var notif models.NewMessageNotification
	if err := json.Unmarshal(data, &notif); err != nil {
		return err
	}
	from := notif.From
	if from == "" {
		from = "someone"
	}
	return c.notifSvc.SendToUserWithSpamCheck(
		notif.UserID,
		"New message",
		"You've received a new message from "+from+".",
		"new-message",
		map[string]interface{}{"from": from},
	)
}

func (c *Consumer) handleConsultationCreated(data []byte) error {
	var notif models.ConsultationCreatedNotification
	if err := json.Unmarshal(data, &notif); err != nil {
		return err
	}

	log.Printf(" Session created: %s between %s and %s", notif.SessionID, notif.ClientName, notif.TherapistName)

	
	if notif.ClientEmail != "" {
		err := c.emailSvc.SendSessionCreatedEmail(
			notif.ClientEmail,
			notif.ClientName,
			"client",
			notif.TherapistName,
			notif.DateStr,
			notif.StartTimeStr,
			notif.Mode,
		)
		if err != nil {
			log.Printf(" Failed to send Gmail to client: %v", err)
		}
	}

	
	if notif.TherapistEmail != "" {
		err := c.emailSvc.SendSessionCreatedEmail(
			notif.TherapistEmail,
			notif.TherapistName,
			"therapist",
			notif.ClientName,
			notif.DateStr,
			notif.StartTimeStr,
			notif.Mode,
		)
		if err != nil {
			log.Printf(" Failed to send Gmail to therapist: %v", err)
		}
	}

	
	clientMsg := fmt.Sprintf("You have an upcoming %s session with %s on %s at %s.",
		strings.ToLower(notif.Mode), notif.TherapistName, notif.DateStr, notif.StartTimeStr)
	_ = c.notifSvc.SendToUser(
		notif.ClientID,
		"Session Confirmed",
		clientMsg,
		"consultation-created",
		map[string]interface{}{
			"sessionId":     notif.SessionID,
			"therapistId":   notif.TherapistID,
			"therapistName": notif.TherapistName,
			"startTime":     notif.StartTimeStr,
			"date":          notif.DateStr,
			"type":          "session_banner",
		},
	)

	
	therapistMsg := fmt.Sprintf("New %s session booked by %s on %s at %s.",
		strings.ToLower(notif.Mode), notif.ClientName, notif.DateStr, notif.StartTimeStr)
	_ = c.notifSvc.SendToUser(
		notif.TherapistID,
		"New Booking",
		therapistMsg,
		"consultation-created",
		map[string]interface{}{
			"sessionId":  notif.SessionID,
			"clientId":   notif.ClientID,
			"clientName": notif.ClientName,
			"startTime":  notif.StartTimeStr,
			"date":       notif.DateStr,
			"type":       "session_banner",
		},
	)

	return nil
}

func (c *Consumer) handleConsultationUpdated(data []byte) error {
	var notif models.ConsultationUpdatedNotification
	if err := json.Unmarshal(data, &notif); err != nil {
		return err
	}
	return c.notifSvc.SendToUser(
		notif.UserID,
		"Consultation updated",
		"Your consultation schedule has been updated.",
		"consultation-updated",
		map[string]interface{}{},
	)
}

func (c *Consumer) handleConsultationReminder(data []byte) error {
	var notif models.ConsultationReminderNotification
	if err := json.Unmarshal(data, &notif); err != nil {
		return err
	}
	return c.notifSvc.SendToUserWithSpamCheck(
		notif.UserID,
		"Upcoming consultation",
		"Reminder: You have a consultation starting soon.",
		"consultation-reminder",
		map[string]interface{}{},
	)
}

func (c *Consumer) handleConsultationCompleted(data []byte) error {
	var notif models.ConsultationCompletedNotification
	if err := json.Unmarshal(data, &notif); err != nil {
		return err
	}
	return c.notifSvc.SendToUser(
		notif.UserID,
		"Consultation completed",
		"Your consultation has been completed successfully.",
		"consultation-completed",
		map[string]interface{}{},
	)
}

func (c *Consumer) handleProfileUpdate(data []byte) error {
	var notif models.ProfileUpdateNotification
	if err := json.Unmarshal(data, &notif); err != nil {
		return err
	}
	return c.notifSvc.SendToUserWithSpamCheck(
		notif.UserID,
		"Profile update required",
		"Your therapist profile needs verification or update.",
		"profile-update",
		map[string]interface{}{},
	)
}

func (c *Consumer) handleProfileApproved(data []byte) error {
	var notif models.ProfileApprovedNotification
	if err := json.Unmarshal(data, &notif); err != nil {
		return err
	}
	return c.notifSvc.SendToUser(
		notif.UserID,
		"Profile approved",
		"Your therapist profile has been approved!",
		"profile-approved",
		map[string]interface{}{},
	)
}

func (c *Consumer) handleProfileRejected(data []byte) error {
	var notif models.ProfileRejectedNotification
	if err := json.Unmarshal(data, &notif); err != nil {
		return err
	}
	return c.notifSvc.SendToUser(
		notif.UserID,
		"Profile rejected",
		"Your therapist profile was rejected. Please review and resubmit.",
		"profile-rejected",
		map[string]interface{}{},
	)
}

func (c *Consumer) handleNewReview(data []byte) error {
	var notif models.NewReviewNotification
	if err := json.Unmarshal(data, &notif); err != nil {
		return err
	}
	reviewer := notif.Reviewer
	if reviewer == "" {
		reviewer = "A client"
	}
	return c.notifSvc.SendToUserWithSpamCheck(
		notif.UserID,
		"New review received",
		reviewer+" just left a review for you.",
		"new-review",
		map[string]interface{}{"reviewer": reviewer},
	)
}

func (c *Consumer) handleNewClient(data []byte) error {
	var notif models.NewClientNotification
	if err := json.Unmarshal(data, &notif); err != nil {
		return err
	}
	clientName := notif.ClientName
	if clientName == "" {
		clientName = "A new user"
	}
	return c.notifSvc.SendToUser(
		notif.UserID,
		"New client connected",
		clientName+" has booked a consultation with you!",
		"new-client",
		map[string]interface{}{"clientName": clientName},
	)
}

func (c *Consumer) handleSystem(data []byte) error {
	var notif models.SystemNotification
	if err := json.Unmarshal(data, &notif); err != nil {
		return err
	}
	return c.notifSvc.SendToUser(
		notif.UserID,
		notif.Title,
		notif.Message,
		"system",
		map[string]interface{}{},
	)
}


func (c *Consumer) handlePostUpvote(data []byte) error {
	var notif models.PostUpvoteNotification
	if err := json.Unmarshal(data, &notif); err != nil {
		return err
	}
	voterName := notif.VoterName
	if voterName == "" {
		voterName = "Someone"
	}
	return c.notifSvc.SendToUser(
		notif.UserID,
		"New Upvote! ⬆️",
		voterName+" upvoted your post.",
		"upvote",
		map[string]interface{}{"postId": notif.PostID},
	)
}

func (c *Consumer) handlePostBookmark(data []byte) error {
	var notif models.PostBookmarkNotification
	if err := json.Unmarshal(data, &notif); err != nil {
		return err
	}
	bookmarkerName := notif.BookmarkerName
	if bookmarkerName == "" {
		bookmarkerName = "Someone"
	}
	return c.notifSvc.SendToUser(
		notif.UserID,
		"Post Bookmarked! 🔖",
		bookmarkerName+" saved your post.",
		"bookmark",
		map[string]interface{}{"postId": notif.PostID},
	)
}

func (c *Consumer) handlePostShare(data []byte) error {
	var notif models.PostShareNotification
	if err := json.Unmarshal(data, &notif); err != nil {
		return err
	}
	sharerName := notif.SharerName
	if sharerName == "" {
		sharerName = "Someone"
	}
	return c.notifSvc.SendToUser(
		notif.UserID,
		"Post Shared! 🔗",
		sharerName+" shared your post.",
		"share",
		map[string]interface{}{"postId": notif.PostID},
	)
}

func (c *Consumer) handlePostComment(data []byte) error {
	var notif models.PostCommentNotification
	if err := json.Unmarshal(data, &notif); err != nil {
		return err
	}
	commenterName := notif.CommenterName
	if commenterName == "" {
		commenterName = "Someone"
	}
	return c.notifSvc.SendToUser(
		notif.UserID,
		"New Comment! 💬",
		commenterName+" commented on your post.",
		"comment",
		map[string]interface{}{
			"postId":    notif.PostID,
			"commentId": notif.CommentID,
		},
	)
}

func (c *Consumer) handleUserFollow(data []byte) error {
	var notif models.UserFollowNotification
	if err := json.Unmarshal(data, &notif); err != nil {
		return err
	}
	followerName := notif.FollowerName
	if followerName == "" {
		followerName = "Someone"
	}
	return c.notifSvc.SendToUser(
		notif.UserID,
		"New Follower! 👤",
		followerName+" started following you.",
		"follow",
		map[string]interface{}{"followerId": notif.FollowerID},
	)
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
