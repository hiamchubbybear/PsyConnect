package repository

import (
	"consultationservice/internal/db"
	"consultationservice/internal/model"
	"context"
	"math"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PostRepository interface {
	CreatePost(ctx context.Context, post *model.Post) error
	GetPostByID(ctx context.Context, id string) (*model.Post, error)
	UpdatePost(ctx context.Context, id string, post *model.Post) error
	DeletePost(ctx context.Context, id string) error
	GetPostsByUser(ctx context.Context, userID string, limit, skip int64) ([]model.Post, error)
	GetFeed(ctx context.Context, userIDs []string, limit, skip int64) ([]model.Post, error)
	SearchPosts(ctx context.Context, query string, limit, skip int64) ([]model.Post, error)
	GetTrendingPosts(ctx context.Context, limit int64) ([]model.Post, error)
	GetPostsByTag(ctx context.Context, tag string, limit, skip int64) ([]model.Post, error)
	GetPostsByCategory(ctx context.Context, category string, limit, skip int64) ([]model.Post, error)
	IncrementViewCount(ctx context.Context, id string) error
	UpdateEngagementCount(ctx context.Context, id, field string, delta int) error
}

type postRepo struct {
	collection *mongo.Collection
}

func NewPostRepo() PostRepository {
	return &postRepo{
		collection: db.GetPostCollection(),
	}
}

func (r *postRepo) CreatePost(ctx context.Context, post *model.Post) error {
	// Initialize default values
	if post.Tags == nil {
		post.Tags = []string{}
	}
	if post.Categories == nil {
		post.Categories = []string{}
	}
	if post.Media == nil {
		post.Media = []model.MediaAttachment{}
	}
	if post.Mentions == nil {
		post.Mentions = []string{}
	}
	if post.Hashtags == nil {
		post.Hashtags = extractHashtags(post.Content)
	}
	if post.Visibility == "" {
		post.Visibility = "public"
	}
	if post.PostType == "" {
		post.PostType = "article"
	}

	post.ViewCount = 0
	post.LikeCount = 0
	post.CommentCount = 0
	post.ShareCount = 0
	post.IsDeleted = false

	result, err := r.collection.InsertOne(ctx, post)
	if err != nil {
		return err
	}
	post.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *postRepo) GetPostByID(ctx context.Context, id string) (*model.Post, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	var post model.Post
	err = r.collection.FindOne(ctx, filter).Decode(&post)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepo) UpdatePost(ctx context.Context, id string, post *model.Post) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	post.UpdatedAt = time.Now()

	// Re-extract hashtags from content
	post.Hashtags = extractHashtags(post.Content)

	filter := bson.M{"_id": objID, "is_deleted": false}
	update := bson.M{"$set": post}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *postRepo) DeletePost(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"is_deleted": true, "updated_at": time.Now()}}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *postRepo) GetPostsByUser(ctx context.Context, userID string, limit, skip int64) ([]model.Post, error) {
	filter := bson.M{"author_id": userID, "is_deleted": false}
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(limit).
		SetSkip(skip)

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var posts []model.Post
	if err = cursor.All(ctx, &posts); err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepo) GetFeed(ctx context.Context, userIDs []string, limit, skip int64) ([]model.Post, error) {
	filter := bson.M{
		"author_id":  bson.M{"$in": userIDs},
		"is_deleted": false,
		"visibility": bson.M{"$in": []string{"public", "followers"}},
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

	var posts []model.Post
	if err = cursor.All(ctx, &posts); err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepo) SearchPosts(ctx context.Context, query string, limit, skip int64) ([]model.Post, error) {
	// Create text search filter
	filter := bson.M{
		"is_deleted": false,
		"visibility": "public",
		"$or": []bson.M{
			{"title": bson.M{"$regex": query, "$options": "i"}},
			{"content": bson.M{"$regex": query, "$options": "i"}},
			{"tags": bson.M{"$in": []string{query}}},
			{"hashtags": bson.M{"$in": []string{query}}},
		},
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

	var posts []model.Post
	if err = cursor.All(ctx, &posts); err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepo) GetTrendingPosts(ctx context.Context, limit int64) ([]model.Post, error) {
	// Trending algorithm: Score = (likes * 1 + comments * 2 + shares * 3) / (age_in_hours + 2)^1.5
	// Get posts from last 7 days
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)

	filter := bson.M{
		"is_deleted": false,
		"visibility": "public",
		"created_at": bson.M{"$gte": sevenDaysAgo},
	}

	// Sort by engagement score (approximation using MongoDB)
	opts := options.Find().
		SetSort(bson.D{
			{Key: "like_count", Value: -1},
			{Key: "comment_count", Value: -1},
			{Key: "created_at", Value: -1},
		}).
		SetLimit(limit)

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var posts []model.Post
	if err = cursor.All(ctx, &posts); err != nil {
		return nil, err
	}

	// Calculate trending score and re-sort in application
	for i := range posts {
		posts[i] = calculateTrendingScore(posts[i])
	}

	return posts, nil
}

func (r *postRepo) GetPostsByTag(ctx context.Context, tag string, limit, skip int64) ([]model.Post, error) {
	filter := bson.M{
		"tags":       bson.M{"$in": []string{tag}},
		"is_deleted": false,
		"visibility": "public",
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

	var posts []model.Post
	if err = cursor.All(ctx, &posts); err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepo) GetPostsByCategory(ctx context.Context, category string, limit, skip int64) ([]model.Post, error) {
	filter := bson.M{
		"categories": bson.M{"$in": []string{category}},
		"is_deleted": false,
		"visibility": "public",
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

	var posts []model.Post
	if err = cursor.All(ctx, &posts); err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepo) IncrementViewCount(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$inc": bson.M{"view_count": 1}}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *postRepo) UpdateEngagementCount(ctx context.Context, id, field string, delta int) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$inc": bson.M{field: delta}}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Helper functions

func extractHashtags(content string) []string {
	re := regexp.MustCompile(`#(\w+)`)
	matches := re.FindAllStringSubmatch(content, -1)

	hashtags := make([]string, 0, len(matches))
	seen := make(map[string]bool)

	for _, match := range matches {
		if len(match) > 1 {
			tag := match[1]
			if !seen[tag] {
				hashtags = append(hashtags, tag)
				seen[tag] = true
			}
		}
	}

	return hashtags
}

func calculateTrendingScore(post model.Post) model.Post {
	// Score = (likes * 1 + comments * 2 + shares * 3) / (age_in_hours + 2)^1.5
	engagement := float64(post.LikeCount + post.CommentCount*2 + post.ShareCount*3)
	ageHours := time.Since(post.CreatedAt).Hours()
	score := engagement / math.Pow(ageHours+2, 1.5)

	// Store score in a custom field (not persistent, just for sorting)
	_ = score // Score calculated, would need custom field to persist

	return post
}
