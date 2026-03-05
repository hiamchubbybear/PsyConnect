





package creds

import (
	"context"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/internal/aws/credentials"
	"go.mongodb.org/mongo-driver/internal/credproviders"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

const (
	
	
	
	
	expiryWindow = 5 * time.Minute
)


type AWSCredentialProvider struct {
	Cred *credentials.Credentials
}


func NewAWSCredentialProvider(httpClient *http.Client, providers ...credentials.Provider) AWSCredentialProvider {
	providers = append(
		providers,
		credproviders.NewEnvProvider(),
		credproviders.NewAssumeRoleProvider(httpClient, expiryWindow),
		credproviders.NewECSProvider(httpClient, expiryWindow),
		credproviders.NewEC2Provider(httpClient, expiryWindow),
	)

	return AWSCredentialProvider{credentials.NewChainCredentials(providers)}
}


func (p AWSCredentialProvider) GetCredentialsDoc(ctx context.Context) (bsoncore.Document, error) {
	creds, err := p.Cred.GetWithContext(ctx)
	if err != nil {
		return nil, err
	}
	builder := bsoncore.NewDocumentBuilder().
		AppendString("accessKeyId", creds.AccessKeyID).
		AppendString("secretAccessKey", creds.SecretAccessKey)
	if token := creds.SessionToken; len(token) > 0 {
		builder.AppendString("sessionToken", token)
	}
	return builder.Build(), nil
}
