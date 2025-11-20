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
)

type ReactionRepository interface {
	AddReaction(ctx context.Context, reaction *model.Reaction) error
	RemoveReaction(ctx context.Context, postID, userID string) error
	GetReactionsByPost(ctx context.Context, postID string) ([]model.Reaction, error)
	GetUserReaction(ctx context.Context, postID, userID string) (*model.Reaction, error)
	CountReactionsByPost(ctx context.Context, postID string) (int64, error)
	CountReactionsByType(ctx context.Context, postID, reactionType string) (int64, error)
}

type reactionRepo struct {
	collection *mongo.Collection
}

func NewReactionRepo() ReactionRepository {
	return &reactionRepo{
		collection: db.GetReactionCollection(),
	}
}

func (r *reactionRepo) AddReaction(ctx context.Context, reaction *model.Reaction) error {
	reaction.CreatedAt = time.Now()

	// Check if user already reacted
	existing, err := r.GetUserReaction(ctx, reaction.PostID, reaction.UserID)
	if err == nil && existing != nil {
		// Update existing reaction
		filter := bson.M{"post_id": reaction.PostID, "user_id": reaction.UserID}
		update := bson.M{"$set": bson.M{"reaction_type": reaction.ReactionType}}
		_, err = r.collection.UpdateOne(ctx, filter, update)
		return err
	}

	result, err := r.collection.InsertOne(ctx, reaction)
	if err != nil {
		return err
	}
	reaction.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *reactionRepo) RemoveReaction(ctx context.Context, postID, userID string) error {
	filter := bson.M{"post_id": postID, "user_id": userID}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("reaction not found")
	}
	return nil
}

func (r *reactionRepo) GetReactionsByPost(ctx context.Context, postID string) ([]model.Reaction, error) {
	filter := bson.M{"post_id": postID}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reactions []model.Reaction
	if err = cursor.All(ctx, &reactions); err != nil {
		return nil, err
	}
	return reactions, nil
}

func (r *reactionRepo) GetUserReaction(ctx context.Context, postID, userID string) (*model.Reaction, error) {
	filter := bson.M{"post_id": postID, "user_id": userID}
	var reaction model.Reaction
	err := r.collection.FindOne(ctx, filter).Decode(&reaction)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &reaction, nil
}

func (r *reactionRepo) CountReactionsByPost(ctx context.Context, postID string) (int64, error) {
	filter := bson.M{"post_id": postID}
	return r.collection.CountDocuments(ctx, filter)
}

func (r *reactionRepo) CountReactionsByType(ctx context.Context, postID, reactionType string) (int64, error) {
	filter := bson.M{"post_id": postID, "reaction_type": reactionType}
	return r.collection.CountDocuments(ctx, filter)
}
