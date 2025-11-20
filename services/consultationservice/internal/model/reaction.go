package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Reaction struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PostID       string             `bson:"post_id" json:"post_id"`
	UserID       string             `bson:"user_id" json:"user_id"`
	ReactionType string             `bson:"reaction_type" json:"reaction_type"` // like, love, laugh, think, sad, angry
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
}

// Reaction types constants
const (
	ReactionLike  = "like"
	ReactionLove  = "love"
	ReactionLaugh = "laugh"
	ReactionThink = "think"
	ReactionSad   = "sad"
	ReactionAngry = "angry"
)
