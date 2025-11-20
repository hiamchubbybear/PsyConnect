package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PostID          string             `bson:"post_id" json:"post_id"`
	UserID          string             `bson:"user_id" json:"user_id"`
	Content         string             `bson:"content" json:"content"`
	ParentCommentID string             `bson:"parent_comment_id,omitempty" json:"parent_comment_id,omitempty"` // For threaded replies

	// Engagement
	LikeCount int `bson:"like_count" json:"like_count"`

	// Control
	IsDeleted bool `bson:"is_deleted" json:"is_deleted"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
