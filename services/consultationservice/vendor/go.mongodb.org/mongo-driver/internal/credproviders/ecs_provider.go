





package credproviders

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/internal/aws/credentials"
)

const (
	
	ecsProviderName = "ECSProvider"

	awsRelativeURI = "http://169.254.170.2/"
)


type ECSProvider struct {
	AwsContainerCredentialsRelativeURIEnv EnvVar

	httpClient *http.Client
	expiration time.Time

	
	
	
	
	
	expiryWindow time.Duration
}


func NewECSProvider(httpClient *http.Client, expiryWindow time.Duration) *ECSProvider {
	return &ECSProvider{
		
		AwsContainerCredentialsRelativeURIEnv: EnvVar("AWS_CONTAINER_CREDENTIALS_RELATIVE_URI"),
		httpClient:                            httpClient,
		expiryWindow:                          expiryWindow,
	}
}


func (e *ECSProvider) RetrieveWithContext(ctx context.Context) (credentials.Value, error) {
	const defaultHTTPTimeout = 10 * time.Second

	v := credentials.Value{ProviderName: ecsProviderName}

	relativeEcsURI := e.AwsContainerCredentialsRelativeURIEnv.Get()
	if len(relativeEcsURI) == 0 {
		return v, errors.New("AWS_CONTAINER_CREDENTIALS_RELATIVE_URI is missing")
	}
	fullURI := awsRelativeURI + relativeEcsURI

	req, err := http.NewRequest(http.MethodGet, fullURI, nil)
	if err != nil {
		return v, err
	}
	req.Header.Set("Accept", "application/json")

	ctx, cancel := context.WithTimeout(ctx, defaultHTTPTimeout)
	defer cancel()
	resp, err := e.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return v, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return v, fmt.Errorf("response failure: %s", resp.Status)
	}

	var ecsResp struct {
		AccessKeyID     string    `json:"AccessKeyId"`
		SecretAccessKey string    `json:"SecretAccessKey"`
		Token           string    `json:"Token"`
		Expiration      time.Time `json:"Expiration"`
	}

	err = json.NewDecoder(resp.Body).Decode(&ecsResp)
	if err != nil {
		return v, err
	}

	v.AccessKeyID = ecsResp.AccessKeyID
	v.SecretAccessKey = ecsResp.SecretAccessKey
	v.SessionToken = ecsResp.Token
	if !v.HasKeys() {
		return v, errors.New("failed to retrieve ECS keys")
	}
	e.expiration = ecsResp.Expiration.Add(-e.expiryWindow)

	return v, nil
}


func (e *ECSProvider) Retrieve() (credentials.Value, error) {
	return e.RetrieveWithContext(context.Background())
}


func (e *ECSProvider) IsExpired() bool {
	return e.expiration.Before(time.Now())
}
