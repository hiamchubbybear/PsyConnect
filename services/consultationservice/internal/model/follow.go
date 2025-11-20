package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Follow struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FollowerID  string             `bson:"follower_id" json:"follower_id"`   // User who follows
	FollowingID string             `bson:"following_id" json:"following_id"` // User being followed
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
}
