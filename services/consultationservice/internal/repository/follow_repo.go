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

type FollowRepository interface {
	FollowUser(ctx context.Context, followerID, followingID string) error
	UnfollowUser(ctx context.Context, followerID, followingID string) error
	GetFollowers(ctx context.Context, userID string, limit, skip int64) ([]model.Follow, error)
	GetFollowing(ctx context.Context, userID string, limit, skip int64) ([]model.Follow, error)
	IsFollowing(ctx context.Context, followerID, followingID string) (bool, error)
	CountFollowers(ctx context.Context, userID string) (int64, error)
	CountFollowing(ctx context.Context, userID string) (int64, error)
}

type followRepo struct {
	collection *mongo.Collection
}

func NewFollowRepo() FollowRepository {
	return &followRepo{
		collection: db.GetFollowCollection(),
	}
}

func (r *followRepo) FollowUser(ctx context.Context, followerID, followingID string) error {
	if followerID == followingID {
		return errors.New("cannot follow yourself")
	}

	// Check if already following
	isFollowing, err := r.IsFollowing(ctx, followerID, followingID)
	if err != nil {
		return err
	}
	if isFollowing {
		return errors.New("already following this user")
	}

	follow := &model.Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
		CreatedAt:   time.Now(),
	}

	result, err := r.collection.InsertOne(ctx, follow)
	if err != nil {
		return err
	}
	follow.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *followRepo) UnfollowUser(ctx context.Context, followerID, followingID string) error {
	filter := bson.M{
		"follower_id":  followerID,
		"following_id": followingID,
	}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("not following this user")
	}
	return nil
}

func (r *followRepo) GetFollowers(ctx context.Context, userID string, limit, skip int64) ([]model.Follow, error) {
	filter := bson.M{"following_id": userID}
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(limit).
		SetSkip(skip)

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var follows []model.Follow
	if err = cursor.All(ctx, &follows); err != nil {
		return nil, err
	}
	return follows, nil
}

func (r *followRepo) GetFollowing(ctx context.Context, userID string, limit, skip int64) ([]model.Follow, error) {
	filter := bson.M{"follower_id": userID}
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(limit).
		SetSkip(skip)

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var follows []model.Follow
	if err = cursor.All(ctx, &follows); err != nil {
		return nil, err
	}
	return follows, nil
}

func (r *followRepo) IsFollowing(ctx context.Context, followerID, followingID string) (bool, error) {
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

func (r *followRepo) CountFollowers(ctx context.Context, userID string) (int64, error) {
	filter := bson.M{"following_id": userID}
	return r.collection.CountDocuments(ctx, filter)
}

func (r *followRepo) CountFollowing(ctx context.Context, userID string) (int64, error) {
	filter := bson.M{"follower_id": userID}
	return r.collection.CountDocuments(ctx, filter)
}
