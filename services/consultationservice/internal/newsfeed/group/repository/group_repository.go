package repository

import (
	"consultationservice/internal/db"
	"consultationservice/internal/newsfeed/group/domain"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GroupRepository interface {
	CreateGroup(ctx context.Context, group *domain.Group) error
	GetGroups(ctx context.Context, category, query string, limit, skip int64) ([]domain.Group, error)
	GetGroupByID(ctx context.Context, id string) (*domain.Group, error)
	AddMember(ctx context.Context, groupID string, userID string) error
	IsMember(ctx context.Context, groupID string, userID string) (bool, error)
	GetMemberCount(ctx context.Context, groupID string) (int, error)
}

type groupRepo struct {
	groupCollection  *mongo.Collection
	memberCollection *mongo.Collection
}

func NewGroupRepo() GroupRepository {
	return &groupRepo{
		groupCollection:  db.GetGroupCollection(),
		memberCollection: db.GetGroupMemberCollection(),
	}
}

func (r *groupRepo) CreateGroup(ctx context.Context, group *domain.Group) error {
	group.CreatedAt = time.Now()
	group.UpdatedAt = time.Now()
	group.MemberCount = 1 

	result, err := r.groupCollection.InsertOne(ctx, group)
	if err != nil {
		return err
	}
	group.ID = result.InsertedID.(primitive.ObjectID)

	
	member := domain.GroupMember{
		GroupID:  group.ID,
		UserID:   group.CreatorID,
		Role:     "admin",
		JoinedAt: time.Now(),
	}
	_, err = r.memberCollection.InsertOne(ctx, member)
	return err
}

func (r *groupRepo) GetGroups(ctx context.Context, category, query string, limit, skip int64) ([]domain.Group, error) {
	filter := bson.M{}
	if category != "" {
		filter["category"] = category
	}
	if query != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": query, "$options": "i"}},
			{"description": bson.M{"$regex": query, "$options": "i"}},
		}
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "member_count", Value: -1}}).
		SetLimit(limit).
		SetSkip(skip)

	cursor, err := r.groupCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var groups []domain.Group
	if err = cursor.All(ctx, &groups); err != nil {
		return nil, err
	}
	return groups, nil
}

func (r *groupRepo) GetGroupByID(ctx context.Context, id string) (*domain.Group, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var group domain.Group
	err = r.groupCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&group)
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *groupRepo) AddMember(ctx context.Context, groupID string, userID string) error {
	gID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return err
	}

	
	isMember, err := r.IsMember(ctx, groupID, userID)
	if err != nil {
		return err
	}
	if isMember {
		return nil 
	}

	member := domain.GroupMember{
		GroupID:  gID,
		UserID:   userID,
		Role:     "member",
		JoinedAt: time.Now(),
	}

	_, err = r.memberCollection.InsertOne(ctx, member)
	if err != nil {
		return err
	}

	
	_, err = r.groupCollection.UpdateOne(
		ctx,
		bson.M{"_id": gID},
		bson.M{"$inc": bson.M{"member_count": 1}},
	)
	return err
}

func (r *groupRepo) IsMember(ctx context.Context, groupID string, userID string) (bool, error) {
	gID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return false, err
	}

	count, err := r.memberCollection.CountDocuments(ctx, bson.M{
		"group_id": gID,
		"user_id":  userID,
	})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *groupRepo) GetMemberCount(ctx context.Context, groupID string) (int, error) {
	gID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return 0, err
	}

	count, err := r.memberCollection.CountDocuments(ctx, bson.M{"group_id": gID})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
