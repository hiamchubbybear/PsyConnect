package repository

import (
	"chatservice/internal/chat/model/model"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ConversationRepository struct {
	collection *mongo.Collection
}

func NewConversationRepository(collection *mongo.Collection) *ConversationRepository {
	return &ConversationRepository{collection: collection}
}

func (r *ConversationRepository) Exists(id string) (bool, error) {
	count, err := r.collection.CountDocuments(context.Background(), bson.M{"_id": id})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *ConversationRepository) CreateConversation(conv *model.Conversation) error {
	_, err := r.collection.InsertOne(context.Background(), conv)
	return err
}
func (r *ConversationRepository) GetConversationByUsers(user1ID, user2ID string) (*model.Conversation, error) {
	filter := bson.M{
		"participants": bson.M{
			"$all": []string{user1ID, user2ID},
		},
	}

	var conv model.Conversation
	err := r.collection.FindOne(context.Background(), filter).Decode(&conv)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	log.Println("Converstation founds", conv)
	return &conv, err
}

func (r *ConversationRepository) GetConversationsByUserID(userID string) ([]*model.Conversation, error) {
	filter := bson.M{
		"participants": bson.M{
			"$elemMatch": bson.M{"$eq": userID},
		},
	}

	cursor, err := r.collection.Find(context.Background(), filter)
	if err != nil {
		log.Printf("[ConversationRepository] Error querying conversations for user %s: %v", userID, err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	var conversations []*model.Conversation
	if err := cursor.All(context.Background(), &conversations); err != nil {
		return nil, err
	}

	log.Printf("[ConversationRepository] Found %d conversations for user %s", len(conversations), userID)
	return conversations, nil
}
