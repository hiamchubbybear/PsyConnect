





package auth

import (
	"context"
	"fmt"
	"io"
	"net/http"

	
	
	
	"crypto/md5"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/operation"
)





const MONGODBCR = "MONGODB-CR"

func newMongoDBCRAuthenticator(cred *Cred, _ *http.Client) (Authenticator, error) {
	source := cred.Source
	if source == "" {
		source = "admin"
	}
	return &MongoDBCRAuthenticator{
		DB:       source,
		Username: cred.Username,
		Password: cred.Password,
	}, nil
}





type MongoDBCRAuthenticator struct {
	DB       string
	Username string
	Password string
}





func (a *MongoDBCRAuthenticator) Auth(ctx context.Context, cfg *Config) error {

	db := a.DB
	if db == "" {
		db = defaultAuthDB
	}

	doc := bsoncore.BuildDocumentFromElements(nil, bsoncore.AppendInt32Element(nil, "getnonce", 1))
	cmd := operation.NewCommand(doc).
		Database(db).
		Deployment(driver.SingleConnectionDeployment{cfg.Connection}).
		ClusterClock(cfg.ClusterClock).
		ServerAPI(cfg.ServerAPI)
	err := cmd.Execute(ctx)
	if err != nil {
		return newError(err, MONGODBCR)
	}
	rdr := cmd.Result()

	var getNonceResult struct {
		Nonce string `bson:"nonce"`
	}

	err = bson.Unmarshal(rdr, &getNonceResult)
	if err != nil {
		return newAuthError("unmarshal error", err)
	}

	doc = bsoncore.BuildDocumentFromElements(nil,
		bsoncore.AppendInt32Element(nil, "authenticate", 1),
		bsoncore.AppendStringElement(nil, "user", a.Username),
		bsoncore.AppendStringElement(nil, "nonce", getNonceResult.Nonce),
		bsoncore.AppendStringElement(nil, "key", a.createKey(getNonceResult.Nonce)),
	)
	cmd = operation.NewCommand(doc).
		Database(db).
		Deployment(driver.SingleConnectionDeployment{cfg.Connection}).
		ClusterClock(cfg.ClusterClock).
		ServerAPI(cfg.ServerAPI)
	err = cmd.Execute(ctx)
	if err != nil {
		return newError(err, MONGODBCR)
	}

	return nil
}


func (a *MongoDBCRAuthenticator) Reauth(_ context.Context, _ *driver.AuthConfig) error {
	return newAuthError("MONGODB-CR does not support reauthentication", nil)
}

func (a *MongoDBCRAuthenticator) createKey(nonce string) string {
	
	
	
	h := md5.New()

	_, _ = io.WriteString(h, nonce)
	_, _ = io.WriteString(h, a.Username)
	_, _ = io.WriteString(h, mongoPasswordDigest(a.Username, a.Password))
	return fmt.Sprintf("%x", h.Sum(nil))
}
