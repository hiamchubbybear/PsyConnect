package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	DB *gorm.DB
}

// Notification model
type Notification struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    string         `gorm:"index;not null" json:"userId"`
	Title     string         `gorm:"not null" json:"title"`
	Body      string         `gorm:"type:text;not null" json:"body"`
	Type      string         `gorm:"default:'system'" json:"type"`
	Metadata  string         `gorm:"type:json" json:"metadata"`
	IsRead    bool           `gorm:"default:false" json:"isRead"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// FCMToken model
type FCMToken struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    string         `gorm:"uniqueIndex;not null" json:"userId"`
	Token     string         `gorm:"not null" json:"token"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func NewDatabase(dsn string) (*Database, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-migrate models
	if err := db.AutoMigrate(&Notification{}, &FCMToken{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("✅ Database connected and migrated")

	return &Database{DB: db}, nil
}

// SaveNotification creates a new notification
func (d *Database) SaveNotification(notification *Notification) error {
	return d.DB.Create(notification).Error
}

// GetNotifications retrieves notifications for a user
func (d *Database) GetNotifications(userID string, limit, offset int) ([]Notification, error) {
	var notifications []Notification
	err := d.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifications).Error
	return notifications, err
}

// MarkAsRead marks a notification as read
func (d *Database) MarkAsRead(notificationID uint, userID string) error {
	return d.DB.Model(&Notification{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Update("is_read", true).Error
}

// SaveFCMToken saves or updates FCM token for a user
func (d *Database) SaveFCMToken(userID, token string) error {
	var fcmToken FCMToken
	result := d.DB.Where("user_id = ?", userID).First(&fcmToken)

	if result.Error == gorm.ErrRecordNotFound {
		// Create new token
		fcmToken = FCMToken{
			UserID: userID,
			Token:  token,
		}
		return d.DB.Create(&fcmToken).Error
	}

	// Update existing token
	fcmToken.Token = token
	return d.DB.Save(&fcmToken).Error
}

// GetFCMToken retrieves FCM token for a user
func (d *Database) GetFCMToken(userID string) (string, error) {
	var fcmToken FCMToken
	err := d.DB.Where("user_id = ?", userID).First(&fcmToken).Error
	if err != nil {
		return "", err
	}
	return fcmToken.Token, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
