





package creds

import (
	"context"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/internal/aws/credentials"
	"go.mongodb.org/mongo-driver/internal/credproviders"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)


type AzureCredentialProvider struct {
	cred *credentials.Credentials
}


func NewAzureCredentialProvider(httpClient *http.Client) AzureCredentialProvider {
	return AzureCredentialProvider{
		credentials.NewCredentials(credproviders.NewAzureProvider(httpClient, 1*time.Minute)),
	}
}


func (p AzureCredentialProvider) GetCredentialsDoc(ctx context.Context) (bsoncore.Document, error) {
	creds, err := p.cred.GetWithContext(ctx)
	if err != nil {
		return nil, err
	}
	builder := bsoncore.NewDocumentBuilder().
		AppendString("accessToken", creds.SessionToken)
	return builder.Build(), nil
}
