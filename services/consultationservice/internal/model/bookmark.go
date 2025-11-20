package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Bookmark struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    string             `bson:"user_id" json:"user_id"`
	PostID    string             `bson:"post_id" json:"post_id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}
