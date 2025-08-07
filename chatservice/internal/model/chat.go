package model

import (
	"log"
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ID             string    `bson:"_id,omitempty"`
	SenderID       string    `json:"Sender" bson:"sender_id"`
	ConversationID string    `json:"conversationId" bson:"conversation_id"`
	Text           string    `json:"text" bson:"text"`
	CreatedAt      time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt      time.Time `json:"updatedAt,omitempty" bson:"updated_at,omitempty"`
}

func NowUTC() time.Time {
	return time.Now().UTC()
}
func NewUUID() string {
	uuid, err := uuid.NewUUID()
	if err != nil {
		log.Printf("Can generate new UUID ")
	}
	return uuid.String()
}
