package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Follow struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FollowerID  string             `bson:"follower_id" json:"follower_id"`
	FollowingID string             `bson:"following_id" json:"following_id"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
}

func NewFollow(followerID, followingID string) *Follow {
	return &Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
		CreatedAt:   time.Now().UTC(),
	}
}

type Bookmark struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    string             `bson:"user_id" json:"user_id"`
	PostID    string             `bson:"post_id" json:"post_id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

func NewBookmark(userID, postID string) *Bookmark {
	return &Bookmark{
		UserID:    userID,
		PostID:    postID,
		CreatedAt: time.Now().UTC(),
	}
}
