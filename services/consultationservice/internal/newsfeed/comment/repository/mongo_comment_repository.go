package repository

import (
	"consultationservice/internal/newsfeed/comment/domain"
	"consultationservice/internal/redis"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoCommentRepository struct {
	collection *mongo.Collection
	redis      redis.RedisStore
}

func NewMongoCommentRepository(collection *mongo.Collection, redis redis.RedisStore) *MongoCommentRepository {
	return &MongoCommentRepository{
		collection: collection,
		redis:      redis,
	}
}

func (r *MongoCommentRepository) CreateComment(ctx context.Context, comment *domain.Comment) error {
	comment.CreatedAt = time.Now().UTC()
	comment.UpdatedAt = time.Now().UTC()
	comment.ReplyCount = 0
	comment.IsDeleted = false

	
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

		
		if comment.Depth > 2 {
			return errors.New("max depth exceeded")
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

	
	if comment.ParentCommentID != "" {
		_ = r.IncrementReplyCount(ctx, comment.ParentCommentID)
	}

	
	r.invalidateCache(ctx, comment.PostID)
	return nil
}

func (r *MongoCommentRepository) UpdateComment(ctx context.Context, id, content string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	update := bson.M{
		"$set": bson.M{
			"content":    content,
			"updated_at": time.Now().UTC(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	
	comment, _ := r.GetCommentByID(ctx, id)
	if comment != nil {
		r.invalidateCache(ctx, comment.PostID)
	}

	return nil
}

func (r *MongoCommentRepository) DeleteComment(ctx context.Context, id string) error {
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

	
	comment, _ := r.GetCommentByID(ctx, id)
	if comment != nil {
		r.invalidateCache(ctx, comment.PostID)
	}

	return nil
}

func (r *MongoCommentRepository) GetCommentByID(ctx context.Context, id string) (*domain.Comment, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	var comment domain.Comment
	err = r.collection.FindOne(ctx, filter).Decode(&comment)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *MongoCommentRepository) GetCommentsByPost(ctx context.Context, postID string, limit, skip int64) ([]domain.Comment, error) {
	filter := bson.M{
		"post_id":           postID,
		"is_deleted":        false,
		"parent_comment_id": bson.M{"$exists": false}, 
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

	var comments []domain.Comment
	if err = cursor.All(ctx, &comments); err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *MongoCommentRepository) GetReplies(ctx context.Context, parentCommentID string) ([]domain.Comment, error) {
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

	var comments []domain.Comment
	if err = cursor.All(ctx, &comments); err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *MongoCommentRepository) GetCommentsTree(ctx context.Context, postID string, maxDepth int, limit, skip int64) ([]domain.Comment, error) {
	
	cacheKey := fmt.Sprintf("post:%s:comments:tree:%d:%d:%d", postID, maxDepth, limit, skip)
	var cachedComments []domain.Comment
	err := r.redis.Get(ctx, cacheKey, &cachedComments)
	if err == nil {
		log.Printf("[CommentRepo] CACHE HIT for key: %s", cacheKey)
		return cachedComments, nil
	}

	log.Printf("[CommentRepo] CACHE MISS for key: %s", cacheKey)

	
	filter := bson.M{
		"post_id":    postID,
		"is_deleted": false,
		"depth":      bson.M{"$lte": maxDepth},
	}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var allComments []domain.Comment
	if err = cursor.All(ctx, &allComments); err != nil {
		return nil, err
	}

	log.Printf("[CommentRepo] Found %d raw comments in DB", len(allComments))

	
	tree := buildCommentTree(allComments)
	log.Printf("[CommentRepo] Built tree with %d root comments", len(tree))

	
	r.redis.Set(ctx, cacheKey, tree)

	return tree, nil
}

func buildCommentTree(comments []domain.Comment) []domain.Comment {
	
	commentMap := make(map[string]*domain.Comment)
	for i := range comments {
		id := comments[i].ID.Hex()
		commentMap[id] = &comments[i]
		comments[i].Replies = []domain.Comment{} 
	}

	
	var rootComments []domain.Comment
	for i := range comments {
		if comments[i].ParentCommentID == "" {
			
			rootComments = append(rootComments, comments[i])
		} else {
			
			if parent, exists := commentMap[comments[i].ParentCommentID]; exists {
				parent.Replies = append(parent.Replies, comments[i])
			}
		}
	}

	return rootComments
}

func (r *MongoCommentRepository) CountCommentsByPost(ctx context.Context, postID string) (int64, error) {
	filter := bson.M{"post_id": postID, "is_deleted": false}
	return r.collection.CountDocuments(ctx, filter)
}

func (r *MongoCommentRepository) IncrementReplyCount(ctx context.Context, commentID string) error {
	objID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$inc": bson.M{"reply_count": 1}}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *MongoCommentRepository) DecrementReplyCount(ctx context.Context, commentID string) error {
	objID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$inc": bson.M{"reply_count": -1}}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *MongoCommentRepository) invalidateCache(ctx context.Context, postID string) {
	
	postCacheKey := fmt.Sprintf("post:%s", postID)
	r.redis.Delete(ctx, postCacheKey)

	
	for _, d := range []int{0, 1, 2, 5} {
		treeCacheKey := fmt.Sprintf("post:%s:comments:tree:%d:20:0", postID, d)
		r.redis.Delete(ctx, treeCacheKey)
	}
}
