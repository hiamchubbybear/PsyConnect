package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PostID          string             `bson:"post_id" json:"post_id"`
	UserID          string             `bson:"user_id" json:"user_id"`
	Content         string             `bson:"content" json:"content" validate:"required,max=500"`
	ParentCommentID string             `bson:"parent_comment_id,omitempty" json:"parent_comment_id,omitempty"`

	
	Depth           int                `bson:"depth" json:"depth"`                     
	Path            string             `bson:"path,omitempty" json:"path,omitempty"`   
	ReplyCount      int                `bson:"reply_count" json:"reply_count"`         

	IsDeleted       bool               `bson:"is_deleted" json:"is_deleted"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`

	
	Replies         []Comment          `bson:"-" json:"replies,omitempty"`             
	AuthorID        string             `bson:"-" json:"author_id,omitempty"`           
	Author          *AuthorInfo        `bson:"-" json:"author,omitempty"`              
}

type AuthorInfo struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar,omitempty"`
}
