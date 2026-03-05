





package creds

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)


type GCPCredentialProvider struct {
	httpClient *http.Client
}


func NewGCPCredentialProvider(httpClient *http.Client) GCPCredentialProvider {
	return GCPCredentialProvider{httpClient}
}


func (p GCPCredentialProvider) GetCredentialsDoc(ctx context.Context) (bsoncore.Document, error) {
	metadataHost := "metadata.google.internal"
	if envhost := os.Getenv("GCE_METADATA_HOST"); envhost != "" {
		metadataHost = envhost
	}
	url := fmt.Sprintf("http://%s/computeMetadata/v1/instance/service-accounts/default/token", metadataHost)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve GCP credentials: %w", err)
	}
	req.Header.Set("Metadata-Flavor", "Google")
	resp, err := p.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve GCP credentials: %w", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve GCP credentials: error reading response body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"unable to retrieve GCP credentials: expected StatusCode 200, got StatusCode: %v. Response body: %s",
			resp.StatusCode,
			body)
	}
	var tokenResponse struct {
		AccessToken string `json:"access_token"`
	}
	
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to retrieve GCP credentials: error reading body JSON: %w (response body: %s)",
			err,
			body)
	}
	if tokenResponse.AccessToken == "" {
		return nil, fmt.Errorf("unable to retrieve GCP credentials: got unexpected empty accessToken from GCP Metadata Server. Response body: %s", body)
	}

	builder := bsoncore.NewDocumentBuilder().AppendString("accessToken", tokenResponse.AccessToken)
	return builder.Build(), nil
}
