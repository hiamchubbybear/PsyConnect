package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title    string             `bson:"title" json:"title"`
	Content  string             `bson:"content" json:"content"`
	AuthorID string             `bson:"author_id" json:"author_id"`

	
	Tags       []string          `bson:"tags" json:"tags"`
	Categories []string          `bson:"categories" json:"categories"`
	Media      []MediaAttachment `bson:"media" json:"media"`
	Mentions   []string          `bson:"mentions" json:"mentions"`     
	Hashtags   []string          `bson:"hashtags" json:"hashtags"`     

	
	Visibility string `bson:"visibility" json:"visibility"` 
	PostType   string `bson:"post_type" json:"post_type"`   

	
	ViewCount     int `bson:"view_count" json:"view_count"`
	UpvoteCount   int `bson:"upvote_count" json:"upvote_count"`
	DownvoteCount int `bson:"downvote_count" json:"downvote_count"`
	CommentCount  int `bson:"comment_count" json:"comment_count"`
	ShareCount    int `bson:"share_count" json:"share_count"`

	
	UserVote     *string `bson:"-" json:"user_vote,omitempty"`     
	UserBookmark bool    `bson:"-" json:"user_bookmark,omitempty"` 

	
	IsDeleted bool `bson:"is_deleted" json:"is_deleted"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type MediaAttachment struct {
	Type    string `bson:"type" json:"type"`                         
	URL     string `bson:"url" json:"url"`
	Caption string `bson:"caption,omitempty" json:"caption,omitempty"`
}
