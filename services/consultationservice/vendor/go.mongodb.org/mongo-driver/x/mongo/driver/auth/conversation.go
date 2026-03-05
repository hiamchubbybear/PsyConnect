





package auth

import (
	"context"

	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)









type SpeculativeConversation interface {
	FirstMessage() (bsoncore.Document, error)
	Finish(ctx context.Context, cfg *Config, firstResponse bsoncore.Document) error
}


type SpeculativeAuthenticator interface {
	CreateSpeculativeConversation() (SpeculativeConversation, error)
}
