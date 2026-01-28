package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title    string             `bson:"title" json:"title"`
	Content  string             `bson:"content" json:"content"`
	AuthorID string             `bson:"author_id" json:"author_id"`

	// Enhanced fields
	Tags       []string          `bson:"tags" json:"tags"`
	Categories []string          `bson:"categories" json:"categories"`
	Media      []MediaAttachment `bson:"media" json:"media"`
	Mentions   []string          `bson:"mentions" json:"mentions"` // @username mentions
	Hashtags   []string          `bson:"hashtags" json:"hashtags"` // #topic hashtags

	// Metadata
	Visibility string `bson:"visibility" json:"visibility"` // public, private, followers
	PostType   string `bson:"post_type" json:"post_type"`   // article, question, discussion, resource

	// Engagement metrics
	ViewCount     int `bson:"view_count" json:"view_count"`
	UpvoteCount   int `bson:"upvote_count" json:"upvote_count"`
	DownvoteCount int `bson:"downvote_count" json:"downvote_count"`
	CommentCount  int `bson:"comment_count" json:"comment_count"`
	ShareCount    int `bson:"share_count" json:"share_count"`

	// User-specific state (not stored in DB, populated at runtime)
	UserVote     *string `bson:"-" json:"user_vote,omitempty"`     // "up", "down", or null
	UserBookmark bool    `bson:"-" json:"user_bookmark,omitempty"` // true if user bookmarked

	// Control
	IsDeleted bool `bson:"is_deleted" json:"is_deleted"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type MediaAttachment struct {
	Type    string `bson:"type" json:"type"` // image, video, document
	URL     string `bson:"url" json:"url"`
	Caption string `bson:"caption,omitempty" json:"caption,omitempty"`
}

func NewPost(title, content, authorID string) *Post {
	return &Post{
		Title:         title,
		Content:       content,
		AuthorID:      authorID,
		Tags:          []string{},
		Categories:    []string{},
		Media:         []MediaAttachment{},
		Mentions:      []string{},
		Hashtags:      []string{},
		Visibility:    "public",
		PostType:      "article",
		ViewCount:     0,
		UpvoteCount:   0,
		DownvoteCount: 0,
		CommentCount:  0,
		ShareCount:    0,
		IsDeleted:     false,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}
}
