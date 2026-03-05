package repository

import (
	"consultationservice/internal/db"
	"consultationservice/internal/model"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BookmarkRepository interface {
	AddBookmark(ctx context.Context, userID, postID string) error
	RemoveBookmark(ctx context.Context, userID, postID string) error
	GetBookmarks(ctx context.Context, userID string, limit, skip int64) ([]model.Bookmark, error)
	IsBookmarked(ctx context.Context, userID, postID string) (bool, error)
	CountBookmarks(ctx context.Context, userID string) (int64, error)
}

type bookmarkRepo struct {
	collection *mongo.Collection
}

func NewBookmarkRepo() BookmarkRepository {
	return &bookmarkRepo{
		collection: db.GetBookmarkCollection(),
	}
}

func (r *bookmarkRepo) AddBookmark(ctx context.Context, userID, postID string) error {
	
	isBookmarked, err := r.IsBookmarked(ctx, userID, postID)
	if err != nil {
		return err
	}
	if isBookmarked {
		return errors.New("post already bookmarked")
	}

	bookmark := &model.Bookmark{
		UserID:    userID,
		PostID:    postID,
		CreatedAt: time.Now(),
	}

	result, err := r.collection.InsertOne(ctx, bookmark)
	if err != nil {
		return err
	}
	bookmark.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *bookmarkRepo) RemoveBookmark(ctx context.Context, userID, postID string) error {
	filter := bson.M{
		"user_id": userID,
		"post_id": postID,
	}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("bookmark not found")
	}
	return nil
}

func (r *bookmarkRepo) GetBookmarks(ctx context.Context, userID string, limit, skip int64) ([]model.Bookmark, error) {
	filter := bson.M{"user_id": userID}
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(limit).
		SetSkip(skip)

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookmarks []model.Bookmark
	if err = cursor.All(ctx, &bookmarks); err != nil {
		return nil, err
	}
	return bookmarks, nil
}

func (r *bookmarkRepo) IsBookmarked(ctx context.Context, userID, postID string) (bool, error) {
	filter := bson.M{
		"user_id": userID,
		"post_id": postID,
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *bookmarkRepo) CountBookmarks(ctx context.Context, userID string) (int64, error) {
	filter := bson.M{"user_id": userID}
	return r.collection.CountDocuments(ctx, filter)
}
