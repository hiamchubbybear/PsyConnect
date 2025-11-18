package model

import "time"

type Conversation struct {
	Id           string    `bson:"_id" json:"id"`
	CreatedAt    time.Time `bson:"created_at" json:"createdAt"`
	Participants []string  `bson:"participants" json:"participants"`
}
