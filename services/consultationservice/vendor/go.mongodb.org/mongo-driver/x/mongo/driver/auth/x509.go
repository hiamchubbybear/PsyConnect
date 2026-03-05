





package auth

import (
	"context"
	"net/http"

	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/operation"
)


const MongoDBX509 = "MONGODB-X509"

func newMongoDBX509Authenticator(cred *Cred, _ *http.Client) (Authenticator, error) {
	
	
	
	return &MongoDBX509Authenticator{User: cred.Username}, nil
}


type MongoDBX509Authenticator struct {
	User string
}

var _ SpeculativeAuthenticator = (*MongoDBX509Authenticator)(nil)



type x509Conversation struct{}

var _ SpeculativeConversation = (*x509Conversation)(nil)


func (c *x509Conversation) FirstMessage() (bsoncore.Document, error) {
	return createFirstX509Message(), nil
}


func createFirstX509Message() bsoncore.Document {
	elements := [][]byte{
		bsoncore.AppendInt32Element(nil, "authenticate", 1),
		bsoncore.AppendStringElement(nil, "mechanism", MongoDBX509),
	}

	return bsoncore.BuildDocument(nil, elements...)
}



func (c *x509Conversation) Finish(context.Context, *Config, bsoncore.Document) error {
	return nil
}


func (a *MongoDBX509Authenticator) CreateSpeculativeConversation() (SpeculativeConversation, error) {
	return &x509Conversation{}, nil
}


func (a *MongoDBX509Authenticator) Auth(ctx context.Context, cfg *Config) error {
	requestDoc := createFirstX509Message()
	authCmd := operation.
		NewCommand(requestDoc).
		Database(sourceExternal).
		Deployment(driver.SingleConnectionDeployment{cfg.Connection}).
		ClusterClock(cfg.ClusterClock).
		ServerAPI(cfg.ServerAPI)
	err := authCmd.Execute(ctx)
	if err != nil {
		return newAuthError("round trip error", err)
	}

	return nil
}


func (a *MongoDBX509Authenticator) Reauth(_ context.Context, _ *driver.AuthConfig) error {
	return newAuthError("X509 does not support reauthentication", nil)
}
