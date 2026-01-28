package repository

import (
	"consultationservice/internal/newsfeed/reaction/domain"
	"consultationservice/internal/redis"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoReactionRepository struct {
	collection *mongo.Collection
	redis      redis.RedisStore
}

func NewMongoReactionRepository(collection *mongo.Collection, redis redis.RedisStore) *MongoReactionRepository {
	return &MongoReactionRepository{
		collection: collection,
		redis:      redis,
	}
}

func (r *MongoReactionRepository) AddReaction(ctx context.Context, reaction *domain.Reaction) error {
	reaction.CreatedAt = time.Now().UTC()

	// Check if user already reacted
	existing, err := r.GetUserReaction(ctx, reaction.PostID, reaction.UserID)
	if err == nil && existing != nil {
		// Update existing reaction
		filter := bson.M{"post_id": reaction.PostID, "user_id": reaction.UserID}
		update := bson.M{"$set": bson.M{"reaction_type": reaction.ReactionType}}
		_, err = r.collection.UpdateOne(ctx, filter, update)

		// Invalidate cache
		r.invalidateCache(ctx, reaction.PostID)
		return err
	}

	result, err := r.collection.InsertOne(ctx, reaction)
	if err != nil {
		return err
	}
	reaction.ID = result.InsertedID.(primitive.ObjectID)

	// Invalidate cache
	r.invalidateCache(ctx, reaction.PostID)
	return nil
}

func (r *MongoReactionRepository) RemoveReaction(ctx context.Context, postID, userID string) error {
	filter := bson.M{"post_id": postID, "user_id": userID}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("reaction not found")
	}

	// Invalidate cache
	r.invalidateCache(ctx, postID)
	return nil
}

func (r *MongoReactionRepository) GetUserReaction(ctx context.Context, postID, userID string) (*domain.Reaction, error) {
	filter := bson.M{"post_id": postID, "user_id": userID}
	var reaction domain.Reaction
	err := r.collection.FindOne(ctx, filter).Decode(&reaction)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &reaction, nil
}

func (r *MongoReactionRepository) GetReactionsByPost(ctx context.Context, postID string) ([]domain.Reaction, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("post:%s:reactions", postID)
	var reactions []domain.Reaction
	err := r.redis.Get(ctx, cacheKey, &reactions)
	if err == nil {
		return reactions, nil
	}

	// Get from DB
	filter := bson.M{"post_id": postID}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &reactions); err != nil {
		return nil, err
	}

	// Cache result
	r.redis.Set(ctx, cacheKey, reactions)
	return reactions, nil
}

func (r *MongoReactionRepository) CountReactionsByType(ctx context.Context, postID, reactionType string) (int64, error) {
	filter := bson.M{"post_id": postID, "reaction_type": reactionType}
	return r.collection.CountDocuments(ctx, filter)
}

func (r *MongoReactionRepository) invalidateCache(ctx context.Context, postID string) {
	cacheKey := fmt.Sprintf("post:%s:reactions", postID)
	r.redis.Delete(ctx, cacheKey)
}
