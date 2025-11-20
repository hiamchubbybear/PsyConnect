package repository

import (
	"consultationservice/internal/db"
	"consultationservice/internal/model"
	"context"
	"time"

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
	comment.LikeCount = 0
	comment.IsDeleted = false

	result, err := r.collection.InsertOne(ctx, comment)
	if err != nil {
		return err
	}
	comment.ID = result.InsertedID.(primitive.ObjectID)
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
