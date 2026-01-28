package repository

import (
	"consultationservice/internal/newsfeed/social/domain"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoFollowRepository struct {
	collection *mongo.Collection
}

func NewMongoFollowRepository(collection *mongo.Collection) *MongoFollowRepository {
	return &MongoFollowRepository{
		collection: collection,
	}
}

func (r *MongoFollowRepository) FollowUser(ctx context.Context, followerID, followingID string) error {
	// Check if already following
	isFollowing, err := r.IsFollowing(ctx, followerID, followingID)
	if err != nil {
		return err
	}
	if isFollowing {
		return errors.New("already following")
	}

	follow := domain.NewFollow(followerID, followingID)
	_, err = r.collection.InsertOne(ctx, follow)
	return err
}

func (r *MongoFollowRepository) UnfollowUser(ctx context.Context, followerID, followingID string) error {
	filter := bson.M{
		"follower_id":  followerID,
		"following_id": followingID,
	}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("not following")
	}
	return nil
}

func (r *MongoFollowRepository) GetFollowers(ctx context.Context, userID string) ([]domain.Follow, error) {
	filter := bson.M{"following_id": userID}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var follows []domain.Follow
	if err = cursor.All(ctx, &follows); err != nil {
		return nil, err
	}
	return follows, nil
}

func (r *MongoFollowRepository) GetFollowing(ctx context.Context, userID string) ([]domain.Follow, error) {
	filter := bson.M{"follower_id": userID}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var follows []domain.Follow
	if err = cursor.All(ctx, &follows); err != nil {
		return nil, err
	}
	return follows, nil
}

func (r *MongoFollowRepository) IsFollowing(ctx context.Context, followerID, followingID string) (bool, error) {
	filter := bson.M{
		"follower_id":  followerID,
		"following_id": followingID,
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *MongoFollowRepository) CountFollowers(ctx context.Context, userID string) (int64, error) {
	filter := bson.M{"following_id": userID}
	return r.collection.CountDocuments(ctx, filter)
}

func (r *MongoFollowRepository) CountFollowing(ctx context.Context, userID string) (int64, error) {
	filter := bson.M{"follower_id": userID}
	return r.collection.CountDocuments(ctx, filter)
}

type MongoBookmarkRepository struct {
	collection *mongo.Collection
}

func NewMongoBookmarkRepository(collection *mongo.Collection) *MongoBookmarkRepository {
	return &MongoBookmarkRepository{
		collection: collection,
	}
}

func (r *MongoBookmarkRepository) CreateBookmark(ctx context.Context, userID, postID string) error {
	// Check if already bookmarked
	isBookmarked, err := r.IsBookmarked(ctx, userID, postID)
	if err != nil {
		return err
	}
	if isBookmarked {
		return errors.New("already bookmarked")
	}

	bookmark := domain.NewBookmark(userID, postID)
	_, err = r.collection.InsertOne(ctx, bookmark)
	return err
}

func (r *MongoBookmarkRepository) DeleteBookmark(ctx context.Context, userID, postID string) error {
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

func (r *MongoBookmarkRepository) GetBookmarks(ctx context.Context, userID string, limit, skip int64) ([]domain.Bookmark, error) {
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

	var bookmarks []domain.Bookmark
	if err = cursor.All(ctx, &bookmarks); err != nil {
		return nil, err
	}
	return bookmarks, nil
}

func (r *MongoBookmarkRepository) IsBookmarked(ctx context.Context, userID, postID string) (bool, error) {
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
