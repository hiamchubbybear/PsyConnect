package repository

import (
	"chatservice/internal/model"
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		log.Println("[ChatRepository] CreateChat error:", err)
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
		log.Println("[ChatRepository] FindChatByID error:", err)
		return nil, errors.New("failed to find chat")
	}
	return &chat, nil
}

func (r *ChatRepository) FindChatsByConversation(conversationID string) ([]model.Chat, error) {
	cursor, err := r.collection.Find(
		context.Background(),
		bson.M{"conversation_id": conversationID},
		options.Find().SetSort(bson.M{"created_at": 1}),
	)
	if err != nil {
		log.Println("[ChatRepository] FindChatsByConversation error:", err)
		return nil, errors.New("failed to find chats")
	}
	defer cursor.Close(context.Background())

	var chats []model.Chat
	if err := cursor.All(context.Background(), &chats); err != nil {
		log.Println("[ChatRepository] Cursor decode error:", err)
		return nil, errors.New("failed to decode chats")
	}

	return chats, nil
}

func (r *ChatRepository) DeleteChatByID(chatID string) (int64, error) {
	res, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": chatID})
	if err != nil {
		log.Println("[ChatRepository] DeleteChatByID error:", err)
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
		log.Println("[ChatRepository] UpdateChatByID error:", err)
		return nil, errors.New("failed to update chat")
	}

	return &updated, nil
}
