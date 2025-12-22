package repository

import (
	"consultationservice/internal/db"
	"consultationservice/internal/model"
	"context"
	"time"
	"log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, comment *model.Comment) error
	UpdateComment(ctx context.Context, id, content string) error
	DeleteComment(ctx context.Context, id string) error
	GetCommentByID(ctx context.Context, id string) (*model.Comment, error)
	GetCommentsByPost(ctx context.Context, postID string, limit, skip int64) ([]model.Comment, error)
	GetReplies(ctx context.Context, parentCommentID string) ([]model.Comment, error)
	CountCommentsByPost(ctx context.Context, postID string) (int64, error)

	// Nested comment support
	GetCommentsTree(ctx context.Context, postID string, maxDepth int, limit, skip int64) ([]model.Comment, error)
	IncrementReplyCount(ctx context.Context, commentID string) error
	DecrementReplyCount(ctx context.Context, commentID string) error
}


type commentRepo struct {
	collection *mongo.Collection
}

func NewCommentRepo() CommentRepository {
	return &commentRepo{
		collection: db.GetCommentCollection(),
	}
}

func (r *commentRepo) CreateComment(ctx context.Context, comment *model.Comment) error {
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()
	comment.ReplyCount = 0
	comment.IsDeleted = false

	// Calculate depth and path for nested comments
	if comment.ParentCommentID != "" {
		parent, err := r.GetCommentByID(ctx, comment.ParentCommentID)
		if err != nil {
			return err
		}

		comment.Depth = parent.Depth + 1
		if parent.Path == "" {
			comment.Path = parent.ID.Hex()
		} else {
			comment.Path = parent.Path + "/" + parent.ID.Hex()
		}

		// Enforce max depth (3 levels: 0, 1, 2)
		if comment.Depth > 2 {
			return mongo.ErrInvalidIndexValue // Use as "max depth exceeded" error
		}
	} else {
		comment.Depth = 0
		comment.Path = ""
	}

	result, err := r.collection.InsertOne(ctx, comment)
	if err != nil {
		return err
	}
	comment.ID = result.InsertedID.(primitive.ObjectID)

	// Increment parent reply count
	if comment.ParentCommentID != "" {
		_ = r.IncrementReplyCount(ctx, comment.ParentCommentID)
	}

	return nil
}

func (r *commentRepo) UpdateComment(ctx context.Context, id, content string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	update := bson.M{
		"$set": bson.M{
			"content":    content,
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *commentRepo) DeleteComment(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"is_deleted": true}}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *commentRepo) GetCommentByID(ctx context.Context, id string) (*model.Comment, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	var comment model.Comment
	err = r.collection.FindOne(ctx, filter).Decode(&comment)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepo) GetCommentsByPost(ctx context.Context, postID string, limit, skip int64) ([]model.Comment, error) {
	filter := bson.M{
		"post_id":            postID,
		"is_deleted":         false,
		"parent_comment_id": bson.M{"$exists": false}, // Only top-level comments
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(limit).
		SetSkip(skip)

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var comments []model.Comment
	if err = cursor.All(ctx, &comments); err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *commentRepo) GetReplies(ctx context.Context, parentCommentID string) ([]model.Comment, error) {
	filter := bson.M{
		"parent_comment_id": parentCommentID,
		"is_deleted":        false,
	}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var comments []model.Comment
	if err = cursor.All(ctx, &comments); err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *commentRepo) CountCommentsByPost(ctx context.Context, postID string) (int64, error) {
	filter := bson.M{"post_id": postID, "is_deleted": false}
	return r.collection.CountDocuments(ctx, filter)
}

func (r *commentRepo) GetCommentsTree(ctx context.Context, postID string, maxDepth int, limit, skip int64) ([]model.Comment, error) {
	// Get all comments for this post (up to maxDepth)
	filter := bson.M{
		"post_id":    postID,
		"is_deleted": false,
		"depth":      bson.M{"$lte": maxDepth},
	}

	log.Printf("[CommentRepo] Fetching comments for PostID: %s (MaxDepth: %d)", postID, maxDepth)
	log.Printf("[CommentRepo] MongoDB Filter: %+v", filter)

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		log.Printf("[CommentRepo] Error during Find: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var allComments []model.Comment
	if err = cursor.All(ctx, &allComments); err != nil {
		log.Printf("[CommentRepo] Error during cursor.All: %v", err)
		return nil, err
	}

	log.Printf("[CommentRepo] Found %d raw comments in DB", len(allComments))

	// Build tree structure
	tree := buildCommentTree(allComments)
	log.Printf("[CommentRepo] Built tree with %d root comments", len(tree))

	return tree, nil
}

func buildCommentTree(comments []model.Comment) []model.Comment {
	// Create map for quick lookup
	commentMap := make(map[string]*model.Comment)
	for i := range comments {
		id := comments[i].ID.Hex()
		commentMap[id] = &comments[i]
		comments[i].Replies = []model.Comment{} // Initialize replies
	}

	// Build tree
	var rootComments []model.Comment
	for i := range comments {
		if comments[i].ParentCommentID == "" {
			// Root comment
			rootComments = append(rootComments, comments[i])
		} else {
			// Child comment - attach to parent
			if parent, exists := commentMap[comments[i].ParentCommentID]; exists {
				parent.Replies = append(parent.Replies, comments[i])
			}
		}
	}

	return rootComments
}

func (r *commentRepo) IncrementReplyCount(ctx context.Context, commentID string) error {
	objID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$inc": bson.M{"reply_count": 1}}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *commentRepo) DecrementReplyCount(ctx context.Context, commentID string) error {
	objID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$inc": bson.M{"reply_count": -1}}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}
