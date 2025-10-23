package repository

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"chatservice/internal/model"
)

type ChatRepository struct {
	collection *mongo.Collection
}

func NewChatRepository(collection *mongo.Collection) *ChatRepository {
	return &ChatRepository{collection: collection}
}

func (r *ChatRepository) CreateChat(chat *model.Chat) (interface{}, error) {
	res, err := r.collection.InsertOne(context.Background(), chat)
	if err != nil {
		log.Println(" CreateChat error:", err)
		return nil, errors.New("failed to insert chat")
	}
	return res.InsertedID, nil
}

func (r *ChatRepository) FindChatByID(chatID string) (*model.Chat, error) {
	var chat model.Chat
	err := r.collection.FindOne(context.Background(), bson.M{"_id": chatID}).Decode(&chat)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		log.Println(" FindChatByID error:", err)
		return nil, errors.New("failed to find chat")
	}
	return &chat, nil
}

func (r *ChatRepository) FindChatsByConversation(conversationID string, limit int, before time.Time) ([]model.Chat, error) {
	filter := bson.M{"conversation_id": conversationID}
	if !before.IsZero() {
		filter["created_at"] = bson.M{"$lt": before}
	}

	opts := options.Find().
		SetSort(bson.M{"created_at": -1}).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var chats []model.Chat
	if err := cursor.All(context.Background(), &chats); err != nil {
		return nil, err
	}
	for i, j := 0, len(chats)-1; i < j; i, j = i+1, j-1 {
		chats[i], chats[j] = chats[j], chats[i]
	}

	return chats, nil
}

func (r *ChatRepository) DeleteChatByID(chatID string) (int64, error) {
	res, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": chatID})
	if err != nil {
		log.Println(" DeleteChatByID error:", err)
		return 0, errors.New("failed to delete chat")
	}
	return res.DeletedCount, nil
}

func (r *ChatRepository) UpdateChatByID(chatID string, newText string) (*model.Chat, error) {
	update := bson.M{
		"$set": bson.M{
			"text":       newText,
			"updated_at": model.NowUTC(),
		},
	}

	var updated model.Chat
	err := r.collection.FindOneAndUpdate(
		context.Background(),
		bson.M{"_id": chatID},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updated)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		log.Println(" UpdateChatByID error:", err)
		return nil, errors.New("failed to update chat")
	}

	return &updated, nil
}
