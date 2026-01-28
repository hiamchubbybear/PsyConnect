package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	VoteUp   = "up"
	VoteDown = "down"
)

type Reaction struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PostID       string             `bson:"post_id" json:"post_id"`
	UserID       string             `bson:"user_id" json:"user_id"`
	ReactionType string             `bson:"reaction_type" json:"reaction_type"` // "up" or "down"
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
}

func NewReaction(postID, userID, reactionType string) *Reaction {
	return &Reaction{
		PostID:       postID,
		UserID:       userID,
		ReactionType: reactionType,
		CreatedAt:    time.Now().UTC(),
	}
}

func (r *Reaction) IsUpvote() bool {
	return r.ReactionType == VoteUp
}

func (r *Reaction) IsDownvote() bool {
	return r.ReactionType == VoteDown
}
