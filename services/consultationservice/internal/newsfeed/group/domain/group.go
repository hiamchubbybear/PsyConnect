package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Group struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Category    string             `bson:"category" json:"category"` // e.g., "Depression", "Anxiety", "General"
	Icon        string             `bson:"icon" json:"icon"`         // URL or icon identifier
	CreatorID   string             `bson:"creator_id" json:"creator_id"`
	MemberCount int                `bson:"member_count" json:"member_count"`
	IsPublic    bool               `bson:"is_public" json:"is_public"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type GroupMember struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	GroupID  primitive.ObjectID `bson:"group_id" json:"group_id"`
	UserID   string             `bson:"user_id" json:"user_id"`
	Role     string             `bson:"role" json:"role"` // "admin", "member"
	JoinedAt time.Time          `bson:"joined_at" json:"joined_at"`
}
