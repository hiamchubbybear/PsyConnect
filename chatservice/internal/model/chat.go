package model

import "time"

type Chat struct {
	ID             string    `json:"id" bson:"_id,omitempty"`
	UserID         string    `json:"userId" bson:"user_id"`
	ConversationID string    `json:"conversationId" bson:"conversation_id"`
	Text           string    `json:"text" bson:"text"`
	CreatedAt      time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt      time.Time `json:"updatedAt,omitempty" bson:"updated_at,omitempty"`
}

func NowUTC() time.Time {
	return time.Now().UTC()
}
